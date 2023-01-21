package app

import (
	"context"
	"time"

	"jwt-enabled-api/internal/auth"
	"jwt-enabled-api/internal/dto"
	"jwt-enabled-api/internal/hasher"
	"jwt-enabled-api/internal/repo"
	"jwt-enabled-api/models"
)

func (a *Application) UserGet(
	ctx context.Context,
	userID int,
) (*models.User, error) {
	user, err := a.repo.UserGet(ctx, userID)
	if err != nil {
		if err == repo.ErrNoRecord {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (a *Application) UserCreate(
	ctx context.Context,
	req *dto.UserCreateRequest,
) (userID int, err error) {
	err = a.repo.Tx(
		ctx,
		nil,
		func(ctx context.Context, tx repo.Service) error {
			// check email have not been used
			_, err := tx.UserGetByEmail(ctx, req.Email)
			if err == nil {
				return ErrUsedEmail
			}
			if err != repo.ErrNoRecord {
				return err
			}

			// hash the request password
			passwordHash, err := a.hasher.GenerateHash([]byte(req.Password))
			if err != nil {
				return err
			}

			// create the user
			insertUser := &models.User{
				Email:        req.Email,
				PasswordHash: string(passwordHash),
			}

			if err := tx.UserCreate(ctx, insertUser); err != nil {
				return err
			}

			// set the user id
			userID = insertUser.ID

			return nil
		},
	)

	return
}

func (app *Application) UserLogin(
	ctx context.Context,
	req *dto.UserLoginRequest,
) (resp *dto.TokenResponse, err error) {
	err = app.repo.Tx(
		ctx,
		nil,
		func(ctx context.Context, tx repo.Service) error {
			// get user by provided email address
			user, err := tx.UserGetByEmail(ctx, req.Email)
			if err != nil {
				if err == repo.ErrNoRecord {
					return ErrNotFound
				}
				return err
			}

			// check provided password matches user password
			err = app.hasher.CompareHash(
				[]byte(user.PasswordHash),
				[]byte(req.Password),
			)
			if err != nil {
				if err == hasher.ErrMismatchedHash {
					return ErrIncorrectPassword
				}
				return err
			}

			// generate token
			jwtToken, jwtTokenExpiresAt, err := app.auth.GenerateJwtToken(
				&auth.Payload{UserID: user.ID},
			)
			if err != nil {
				return err
			}
			refreshToken, refreshTokenExpiresAt, err := app.auth.GenerateRefreshToken()
			if err != nil {
				return err
			}

			// hash and then save the refresh token
			refreshTokenHash, err := app.hasher.GenerateHash(
				[]byte(refreshToken),
			)
			if err != nil {
				return err
			}

			err = tx.TokenCreate(ctx, &models.Token{
				TokenHash: string(refreshTokenHash),
				UserID:    user.ID,
				ExpiresAt: refreshTokenExpiresAt,
			})
			if err != nil {
				return err
			}

			// set response
			resp = &dto.TokenResponse{
				JwtTokenResponse: dto.JwtTokenResponse{
					JwtToken:     jwtToken,
					JwtExpiresAt: jwtTokenExpiresAt.Unix(),
				},
				RefreshToken:     refreshToken,
				RefreshExpiresAt: refreshTokenExpiresAt.Unix(),
				UserID:           user.ID,
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (app *Application) UserLogout(
	ctx context.Context,
	userID int,
	refreshToken string,
) error {
	err := app.repo.Tx(
		ctx,
		nil,
		func(ctx context.Context, tx repo.Service) error {
			// get token
			token, err := tx.TokenGet(ctx, userID, refreshToken)
			if err != nil {
				if err == repo.ErrNoRecord {
					return ErrNotFound
				}
				return err
			}
			// invalidate token
			err = tx.TokenUpdate(ctx, token.ID, map[string]any{
				models.TokenColumns.ExpiresAt: time.Now(),
			})
			if err != nil {
				// repo.ErrNoRecord have been checked before at top?
				return err
			}
			return nil
		},
	)
	return err
}

func (app *Application) UserRefreshToken(
	ctx context.Context,
	userID int,
	refreshToken string,
) (resp *dto.JwtTokenResponse, err error) {
	// check token exists
	token, err := app.repo.TokenGet(ctx, userID, refreshToken)
	if err != nil {
		if err == repo.ErrNoRecord {
			return nil, ErrNotFound
		}
		return nil, err
	}
	// create the new jwt token
	jwtToken, expiresAt, err := app.auth.GenerateJwtToken(
		&auth.Payload{UserID: token.UserID},
	)
	if err != nil {
		return nil, err
	}
	// return jwt token
	return &dto.JwtTokenResponse{
		JwtToken:     jwtToken,
		JwtExpiresAt: expiresAt.Unix(),
	}, nil
}
