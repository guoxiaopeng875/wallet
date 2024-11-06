package transaction

import "context"

// Repository defines the repository for transaction.
type Repository interface {
	ListByWalletID(ctx context.Context, walletID uint) ([]Transaction, error)
	Create(ctx context.Context, transaction *Transaction) error
}
