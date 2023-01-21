package repo

import (
	"context"
	"database/sql"

	"jwt-enabled-api/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (repo *Repository) UserGet(
	ctx context.Context,
	userID int,
) (*models.User, error) {
	user, err := models.Users(
		models.UserWhere.ID.EQ(userID),
	).One(ctx, repo.exec)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return user, nil
}

func (repo *Repository) UserGetByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	user, err := models.Users(
		models.UserWhere.Email.EQ(email),
	).One(ctx, repo.exec)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return user, nil
}

func (repo *Repository) UserCreate(
	ctx context.Context,
	user *models.User,
) error {
	return user.Insert(ctx, repo.exec, boil.Infer())
}
