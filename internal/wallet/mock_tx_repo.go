package wallet

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/wallet/transaction"
)

type MockTransactionRepository struct {
	transactions []transaction.Transaction
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		transactions: make([]transaction.Transaction, 0),
	}
}

func (m *MockTransactionRepository) Create(ctx context.Context, tx *transaction.Transaction) error {
	m.transactions = append(m.transactions, *tx)
	return nil
}

func (m *MockTransactionRepository) ListByWalletID(ctx context.Context, walletID uint) ([]transaction.Transaction, error) {
	result := make([]transaction.Transaction, 0)
	for _, tx := range m.transactions {
		if tx.FromWalletID == walletID || tx.ToWalletID == walletID {
			result = append(result, tx)
		}
	}
	return result, nil
}
