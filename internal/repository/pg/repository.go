package pg

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/wallet"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func NewConnect(ctx context.Context, dsn string) (*pgx.Conn, func(), error) {
	closer := func() {}
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, closer, err
	}
	closer = func() {
		_ = conn.Close(ctx)
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, closer, err
	}
	return conn, closer, nil
}

func NewDBTx(repo *Repository) wallet.DBTx {
	return repo
}

type contextTxKey struct{}

func (repo *Repository) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return wrapError(repo.execTx(ctx, fn))
}
func (repo *Repository) execTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	ctx = context.WithValue(ctx, contextTxKey{}, tx)
	if err := fn(ctx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (repo *Repository) DB(ctx context.Context) *pgx.Conn {
	tx, ok := ctx.Value(contextTxKey{}).(*pgx.Conn)
	if ok {
		return tx
	}
	return repo.db
}
