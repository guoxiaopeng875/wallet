package pg

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/wallet/transaction"
	"github.com/jackc/pgx/v5"
)

type transactionRepository struct {
	*Repository
}

func NewTransactionRepository(repo *Repository) transaction.Repository {
	return &transactionRepository{repo}
}

func (t *transactionRepository) ListByWalletID(ctx context.Context, walletID uint) ([]transaction.Transaction, error) {
	list, err := t.listByWalletID(ctx, walletID)
	return list, wrapError(err)
}
func (t *transactionRepository) listByWalletID(ctx context.Context, walletID uint) ([]transaction.Transaction, error) {
	rows, err := t.DB(ctx).Query(ctx, "select * from transactions where from_wallet_id = $1 or to_wallet_id = $2", walletID, walletID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[transaction.Transaction])
}

func (t *transactionRepository) Create(ctx context.Context, transaction *transaction.Transaction) error {
	_, err := t.DB(ctx).Exec(
		ctx,
		"insert into transactions (method, tx_at, amount, from_wallet_id, to_wallet_id) values ($1, $2, $3, $4, $5 )",
		transaction.Method, transaction.TxAt, transaction.Amount, transaction.FromWalletID, transaction.ToWalletID,
	)
	return err
}
