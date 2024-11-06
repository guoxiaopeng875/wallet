package mocks

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/wallet"
	"github.com/guoxiaopeng875/wallet/internal/wallet/transaction"
	"github.com/shopspring/decimal"
)

type MockUseCase struct {
	OnDeposit            func(ctx context.Context, walletID uint, amount decimal.Decimal) error
	OnWithdraw           func(ctx context.Context, walletID uint, amount decimal.Decimal) error
	OnTransfer           func(ctx context.Context, fromID, toID uint, amount decimal.Decimal) error
	OnWallet             func(ctx context.Context, walletID uint) (*wallet.Wallet, error)
	OnWalletTransactions func(ctx context.Context, walletID uint) ([]transaction.Transaction, error)
}

func (m *MockUseCase) Deposit(ctx context.Context, walletID uint, amount decimal.Decimal) error {
	return m.OnDeposit(ctx, walletID, amount)
}

func (m *MockUseCase) Withdraw(ctx context.Context, walletID uint, amount decimal.Decimal) error {
	return m.OnWithdraw(ctx, walletID, amount)
}

func (m *MockUseCase) Transfer(ctx context.Context, fromID, toID uint, amount decimal.Decimal) error {
	return m.OnTransfer(ctx, fromID, toID, amount)
}

func (m *MockUseCase) Wallet(ctx context.Context, walletID uint) (*wallet.Wallet, error) {
	return m.OnWallet(ctx, walletID)
}

func (m *MockUseCase) WalletTransactions(ctx context.Context, walletID uint) ([]transaction.Transaction, error) {
	return m.OnWalletTransactions(ctx, walletID)
}
