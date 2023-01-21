package repo

import (
	"context"
	"database/sql"
)

func (repo *Repository) Tx(
	ctx context.Context,
	opts *sql.TxOptions,
	fn func(context.Context, Service) error,
) (err error) {
	var tx *sql.Tx
	tx, err = repo.exec.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	r := NewRepository(noOpBeginTx{tx})
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = fn(ctx, r)
	return
}

type noOpBeginTx struct {
	*sql.Tx
}

var _ Executor = noOpBeginTx{}

func (tx noOpBeginTx) BeginTx(
	context.Context,
	*sql.TxOptions,
) (*sql.Tx, error) {
	// intended no operation:
	// it is here just to satisfy type system to allow '*sql.Tx' as an 'Executor' to pass to 'NewRepository'.
	// but it is totally safe as the 'ServiceTx' objest is down-scaled to a 'Service' object that have no 'BeginTx' method.
	return nil, nil
}
