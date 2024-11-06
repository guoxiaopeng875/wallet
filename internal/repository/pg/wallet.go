package pg

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/pkg/errors"
	"github.com/guoxiaopeng875/wallet/internal/wallet"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type walletRepository struct {
	*Repository
}

func NewWalletRepository(repo *Repository) wallet.Repository {
	return &walletRepository{repo}
}

func (wp *walletRepository) Get(ctx context.Context, id uint) (*wallet.Wallet, error) {
	var w wallet.Wallet
	if err := wp.DB(ctx).QueryRow(ctx, "select id, balance from wallets where id = $1", id).Scan(&w.ID, &w.Balance); err != nil {
		return nil, wrapError(err)
	}
	return &w, nil
}

func (wp *walletRepository) UpdateBalance(ctx context.Context, wallet *wallet.Wallet, amount decimal.Decimal) error {
	ct, err := wp.DB(ctx).Exec(ctx, "update wallets set balance = balance + $1 where id = $2 and balance = $3", amount, wallet.ID, wallet.Balance)
	if err != nil {
		return err
	}
	if ct.RowsAffected() != 1 {
		logrus.Warnf("wallet %d balance update failed, oldBalance=%v, amount=%s", wallet.ID, wallet.Balance, amount)
		return errors.RecordNotFound
	}

	return nil
}
