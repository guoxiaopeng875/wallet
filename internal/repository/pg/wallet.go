package pg

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/wallet"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type walletRepository struct {
	db *pgx.Conn
}

func (wp *walletRepository) Get(ctx context.Context, id uint) (*wallet.Wallet, error) {
	var w wallet.Wallet
	if err := wp.db.QueryRow(ctx, "select id, balance from wallets where id = $1", id).Scan(&w); err != nil {
		return nil, wrapError(err)
	}
	return &w, nil
}

func (wp *walletRepository) UpdateBalance(ctx context.Context, wallet *wallet.Wallet, amount decimal.Decimal) error {
	//TODO implement me
	panic("implement me")
}
