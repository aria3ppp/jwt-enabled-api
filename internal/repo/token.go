package repo

import (
	"context"
	"database/sql"

	"jwt-enabled-api/models"

	"github.com/blockloop/scan"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (repo *Repository) TokenGet(
	ctx context.Context,
	userID int,
	refreshToken string,
) (token *models.Token, err error) {
	// token = new(models.Token)

	// // err = repo.exec.QueryRowContext(ctx, tokenGetQuery, refreshToken).
	// // 	Scan(&token.ID, &token.TokenHash, &token.UserID, &token.ExpiresAt)
	// // if err != nil {
	// // 	if err == sql.ErrNoRows {
	// // 		return nil, ErrNoRecord
	// // 	}
	// // 	return nil, err
	// // }

	// rows, err := repo.exec.QueryContext(ctx, tokenGetQuery, refreshToken)
	// if err != nil {
	// 	return nil, err
	// }
	// err = scan.RowStrict(token, rows)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, ErrNoRecord
	// 	}
	// 	return nil, err
	// }

	// return token, nil

	token = new(models.Token)

	// err = repo.exec.QueryRowContext(ctx, tokenGetQuery, userID, refreshToken).
	// 	Scan(&token.ID, &token.TokenHash, &token.UserID, &token.ExpiresAt)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, ErrNoRecord
	// 	}
	// 	return nil, err
	// }

	rows, err := repo.exec.QueryContext(
		ctx,
		tokenGetQuery,
		userID,
		refreshToken,
	)
	if err != nil {
		return nil, err
	}
	err = scan.RowStrict(token, rows)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return token, nil
}

func (repo *Repository) TokenCreate(
	ctx context.Context,
	token *models.Token,
) error {
	return token.Insert(ctx, repo.exec, boil.Infer())
}

func (repo *Repository) TokenUpdate(
	ctx context.Context,
	tokenID int,
	cols map[string]any,
) error {
	rowsAff, err := models.Tokens(
		models.TokenWhere.ID.EQ(tokenID),
	).UpdateAll(ctx, repo.exec, cols)
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return ErrNoRecord
	}
	return nil
}
