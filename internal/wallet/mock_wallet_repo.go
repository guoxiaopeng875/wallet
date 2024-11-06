package wallet

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/pkg/errors"
	"github.com/shopspring/decimal"
)

type MockRepository struct {
	wallets map[uint]*Wallet
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		wallets: make(map[uint]*Wallet),
	}
}

func (m *MockRepository) Get(ctx context.Context, id uint) (*Wallet, error) {
	if w, exists := m.wallets[id]; exists {
		return w, nil
	}
	return nil, errors.RecordNotFound
}

func (m *MockRepository) UpdateBalance(ctx context.Context, w *Wallet, amount decimal.Decimal) error {
	if _, exists := m.wallets[w.ID]; !exists {
		return errors.RecordNotFound
	}
	w.Balance = w.Balance.Add(amount)
	m.wallets[w.ID] = w
	return nil
}

func (m *MockRepository) AddWallet(w *Wallet) {
	m.wallets[w.ID] = w
}
