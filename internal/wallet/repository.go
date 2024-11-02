package wallet

import (
	"context"
	"github.com/shopspring/decimal"
)

// Repository defines the repository for wallet.
type Repository interface {
	// Get gets the wallet by id.
	Get(ctx context.Context, id uint) (*Wallet, error)
	// UpdateBalance updates the balance of the wallet.
	UpdateBalance(ctx context.Context, wallet *Wallet, amount decimal.Decimal) error
}
