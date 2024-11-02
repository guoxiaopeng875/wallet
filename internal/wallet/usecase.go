package wallet

import (
	"context"
	"fmt"
	"github.com/guoxiaopeng875/wallet/internal/pkg/errors"
	transaction2 "github.com/guoxiaopeng875/wallet/internal/wallet/transaction"
	"time"

	"github.com/shopspring/decimal"
)

// UseCase defines use cases for the wallet.
type UseCase interface {
	// Deposit adds the specified amount to the wallet balance.
	// Returns an error if the amount is not positive or if the wallet doesn't exist.
	Deposit(ctx context.Context, walletID uint, amount decimal.Decimal) error

	// Withdraw subtracts the specified amount from the wallet balance.
	// Returns an error if the amount is not positive, if the wallet doesn't exist,
	// or if the wallet has insufficient funds.
	Withdraw(ctx context.Context, walletID uint, amount decimal.Decimal) error

	// Transfer sends the specified amount from one wallet to another.
	// Returns an error if:
	// - The amount is not positive
	// - Either wallet doesn't exist
	// - The source wallet has insufficient funds
	// - There's a concurrent modification conflict
	Transfer(ctx context.Context, fromWalletID, toWalletID uint, amount decimal.Decimal) error

	// Wallet retrieves wallet information by its ID.
	// Returns the wallet details or an error if the wallet doesn't exist.
	Wallet(ctx context.Context, walletID uint) (*Wallet, error)

	// WalletTransactions retrieves all transactions associated with the specified wallet.
	// Returns a list of transactions or an error if the wallet doesn't exist.
	WalletTransactions(ctx context.Context, walletID uint) ([]*transaction2.Transaction, error)
}

// Locker .
type Locker interface {
	Lock(ctx context.Context, key string) error
	Unlock(ctx context.Context, key string)
}

// DBTx is database transaction.
type DBTx interface {
	ExecTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// useCase implements UseCase.
type useCase struct {
	repo   Repository
	txRepo transaction2.Repository
	dbTx   DBTx
	locker Locker
}

func (u *useCase) Deposit(ctx context.Context, walletID uint, amount decimal.Decimal) error {
	if !amount.IsPositive() {
		return errors.InvalidArgs.WithCause(fmt.Errorf("deposit amount must be positive: %v", amount))
	}
	wallet, err := u.repo.Get(ctx, walletID)
	if err != nil {
		return err
	}
	return u.dbTx.ExecTx(ctx, func(ctx context.Context) error {
		if err := u.repo.UpdateBalance(ctx, wallet, amount); err != nil {
			return err
		}
		return u.txRepo.Create(ctx, &transaction2.Transaction{
			Method:     transaction2.MethodDeposit,
			TxAt:       time.Now(),
			Amount:     amount,
			ToWalletID: wallet.ID,
		})
	})
}

func (u *useCase) Withdraw(ctx context.Context, walletID uint, amount decimal.Decimal) error {
	if !amount.IsPositive() {
		return errors.InvalidArgs.WithCause(fmt.Errorf("withdraw amount must be positive: %v", amount))
	}
	wallet, err := u.repo.Get(ctx, walletID)
	if err != nil {
		return err
	}
	if err := wallet.CheckBalance(amount); err != nil {
		return err
	}
	return u.dbTx.ExecTx(ctx, func(ctx context.Context) error {
		if err := u.repo.UpdateBalance(ctx, wallet, amount.Neg()); err != nil {
			return err
		}
		return u.txRepo.Create(ctx, &transaction2.Transaction{
			Method:       transaction2.MethodWithdraw,
			TxAt:         time.Now(),
			Amount:       amount,
			FromWalletID: wallet.ID,
		})
	})
}

func (u *useCase) Transfer(ctx context.Context, fromWalletID, toWalletID uint, amount decimal.Decimal) error {
	if !amount.IsPositive() {
		return errors.InvalidArgs.WithCause(fmt.Errorf("tranfer amount must be positive: %v", amount))
	}
	fromWallet, err := u.repo.Get(ctx, fromWalletID)
	if err != nil {
		return err
	}
	if err := fromWallet.CheckBalance(amount); err != nil {
		return err
	}
	toWallet, err := u.repo.Get(ctx, toWalletID)
	if err != nil {
		return err
	}
	if err := toWallet.CheckBalance(amount); err != nil {
		return err
	}

	return u.dbTx.ExecTx(ctx, func(ctx context.Context) error {
		if err := u.repo.UpdateBalance(ctx, fromWallet, amount.Neg()); err != nil {
			return err
		}
		if err := u.repo.UpdateBalance(ctx, toWallet, amount); err != nil {
			return err
		}
		return u.txRepo.Create(ctx, &transaction2.Transaction{
			Method:       transaction2.MethodTransfer,
			TxAt:         time.Now(),
			Amount:       amount,
			FromWalletID: fromWallet.ID,
			ToWalletID:   toWallet.ID,
		})
	})
}

func (u *useCase) Wallet(ctx context.Context, walletID uint) (*Wallet, error) {
	return u.repo.Get(ctx, walletID)
}

func (u *useCase) WalletTransactions(ctx context.Context, walletID uint) ([]*transaction2.Transaction, error) {
	wallet, err := u.repo.Get(ctx, walletID)
	if err != nil {
		return nil, err
	}
	return u.txRepo.ListByWalletID(ctx, wallet.ID)
}
