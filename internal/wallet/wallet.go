package wallet

import (
	"github.com/guoxiaopeng875/wallet/internal/pkg/errors"
	"github.com/shopspring/decimal"
)

// Wallet defines the wallet entity
type Wallet struct {
	ID      uint            `json:"id"`
	Balance decimal.Decimal `json:"balance"`
}

// CheckBalance checks if the wallet has enough balance
func (w *Wallet) CheckBalance(amount decimal.Decimal) error {
	if w.Balance.LessThan(amount) {
		return errors.InsufficientBalance
	}
	return nil
}

func (w *Wallet) Deposit(amount decimal.Decimal) {
	w.Balance = w.Balance.Add(amount)
}

func (w *Wallet) Withdraw(amount decimal.Decimal) {
	w.Balance = w.Balance.Sub(amount)
}

func (w *Wallet) Transfer(toWallet *Wallet, amount decimal.Decimal) {
	w.Withdraw(amount)
	toWallet.Deposit(amount)
}
