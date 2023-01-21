package repo

import (
	"context"
	"database/sql"

	"jwt-enabled-api/models"
)

type ServiceTx interface {
	Tx(
		context.Context,
		*sql.TxOptions,
		func(context.Context, Service) error,
	) error
	Service
}

type Service interface {
	UserGet(
		ctx context.Context,
		userID int,
	) (*models.User, error)
	UserGetByEmail(
		ctx context.Context,
		email string,
	) (*models.User, error)
	UserCreate(
		ctx context.Context,
		user *models.User,
	) error

	TokenGet(
		ctx context.Context,
		userID int,
		refreshToken string,
	) (token *models.Token, err error)
	TokenCreate(
		ctx context.Context,
		token *models.Token,
	) error
	TokenUpdate(
		ctx context.Context,
		tokenID int,
		cols map[string]any,
	) error
}

type Repository struct {
	exec Executor
}

var _ ServiceTx = (*Repository)(nil)

type Executor interface {
	Exec(string, ...any) (sql.Result, error)
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

func NewRepository(exec Executor) *Repository {
	return &Repository{exec}
}
