package wallet

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/wallet/transaction"
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

type mockDBTx struct{}

func (m *mockDBTx) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func setupTest(t *testing.T) (UseCase, *MockRepository, *MockTransactionRepository) {
	repo := NewMockRepository()
	txRepo := NewMockTransactionRepository()
	dbTx := &mockDBTx{}
	uc := NewUseCase(repo, txRepo, dbTx)

	// Add test wallets
	repo.AddWallet(&Wallet{ID: 1, Balance: decimal.NewFromFloat(1000)})
	repo.AddWallet(&Wallet{ID: 2, Balance: decimal.NewFromFloat(500)})

	return uc, repo, txRepo
}

func TestUseCase_Deposit(t *testing.T) {
	tests := []struct {
		name      string
		walletID  uint
		amount    decimal.Decimal
		wantErr   bool
		setupFunc func(*MockRepository)
	}{
		{
			name:     "successful deposit",
			walletID: 1,
			amount:   decimal.NewFromFloat(100),
			wantErr:  false,
		},
		{
			name:     "negative amount",
			walletID: 1,
			amount:   decimal.NewFromFloat(-100),
			wantErr:  true,
		},
		{
			name:     "wallet not found",
			walletID: 999,
			amount:   decimal.NewFromFloat(100),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo, _ := setupTest(t)
			if tt.setupFunc != nil {
				tt.setupFunc(repo)
			}

			err := uc.Deposit(context.Background(), tt.walletID, tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deposit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_Withdraw(t *testing.T) {
	tests := []struct {
		name     string
		walletID uint
		amount   decimal.Decimal
		wantErr  bool
	}{
		{
			name:     "successful withdraw",
			walletID: 1,
			amount:   decimal.NewFromFloat(100),
			wantErr:  false,
		},
		{
			name:     "insufficient balance",
			walletID: 1,
			amount:   decimal.NewFromFloat(2000),
			wantErr:  true,
		},
		{
			name:     "negative amount",
			walletID: 1,
			amount:   decimal.NewFromFloat(-100),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, _, _ := setupTest(t)
			err := uc.Withdraw(context.Background(), tt.walletID, tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Withdraw() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_Transfer(t *testing.T) {
	tests := []struct {
		name         string
		fromWalletID uint
		toWalletID   uint
		amount       decimal.Decimal
		wantErr      bool
	}{
		{
			name:         "successful transfer",
			fromWalletID: 1,
			toWalletID:   2,
			amount:       decimal.NewFromFloat(100),
			wantErr:      false,
		},
		{
			name:         "insufficient balance",
			fromWalletID: 1,
			toWalletID:   2,
			amount:       decimal.NewFromFloat(2000),
			wantErr:      true,
		},
		{
			name:         "negative amount",
			fromWalletID: 1,
			toWalletID:   2,
			amount:       decimal.NewFromFloat(-100),
			wantErr:      true,
		},
		{
			name:         "wallet not found",
			fromWalletID: 999,
			toWalletID:   2,
			amount:       decimal.NewFromFloat(100),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, _, _ := setupTest(t)
			err := uc.Transfer(context.Background(), tt.fromWalletID, tt.toWalletID, tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_WalletTransactions(t *testing.T) {
	ctx := context.Background()
	uc, _, txRepo := setupTest(t)

	// Create some test transactions
	txs := []struct {
		method       string
		fromWalletID uint
		toWalletID   uint
		amount       float64
	}{
		{"deposit", 0, 1, 100},
		{"withdraw", 1, 0, 50},
		{"transfer", 1, 2, 30},
	}

	for _, tx := range txs {
		_ = txRepo.Create(ctx, &transaction.Transaction{
			Method:       transaction.Method(tx.method),
			TxAt:         time.Now(),
			Amount:       decimal.NewFromFloat(tx.amount),
			FromWalletID: tx.fromWalletID,
			ToWalletID:   tx.toWalletID,
		})
	}

	// Test retrieving transactions
	transactions, err := uc.WalletTransactions(ctx, 1)
	if err != nil {
		t.Errorf("WalletTransactions() error = %v", err)
	}

	if len(transactions) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(transactions))
	}
}
