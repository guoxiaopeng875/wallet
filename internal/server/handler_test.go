package server

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/guoxiaopeng875/wallet/internal/pkg/errors"
	"github.com/guoxiaopeng875/wallet/internal/server/mocks"
	"github.com/guoxiaopeng875/wallet/internal/wallet"
	"github.com/guoxiaopeng875/wallet/internal/wallet/transaction"
	"github.com/shopspring/decimal"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_Deposit(t *testing.T) {
	tests := []struct {
		name       string
		walletID   string
		reqBody    interface{}
		mockSetup  func(*mocks.MockUseCase)
		wantStatus int
	}{
		{
			name:     "successful deposit",
			walletID: "1",
			reqBody: DepositRequest{
				Amount: decimal.NewFromFloat(100),
			},
			mockSetup: func(m *mocks.MockUseCase) {
				m.OnDeposit = func(ctx context.Context, id uint, amount decimal.Decimal) error {
					return nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:     "invalid wallet id",
			walletID: "invalid",
			reqBody: DepositRequest{
				Amount: decimal.NewFromFloat(100),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid request body",
			walletID:   "1",
			reqBody:    "invalid json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "usecase error",
			walletID: "1",
			reqBody: DepositRequest{
				Amount: decimal.NewFromFloat(100),
			},
			mockSetup: func(m *mocks.MockUseCase) {
				m.OnDeposit = func(ctx context.Context, id uint, amount decimal.Decimal) error {
					return errors.InternalServer
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mocks.MockUseCase{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockUC)
			}

			h := NewHandler(mockUC)
			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/wallets/"+tt.walletID+"/deposit", bytes.NewReader(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.walletID})
			w := httptest.NewRecorder()

			h.Deposit(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Deposit() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_Withdraw(t *testing.T) {
	tests := []struct {
		name       string
		walletID   string
		reqBody    interface{}
		setupMock  func(*mocks.MockUseCase)
		wantStatus int
	}{
		{
			name:     "successful withdrawal",
			walletID: "1",
			reqBody: WithdrawRequest{
				Amount: decimal.NewFromFloat(50.0),
			},
			setupMock: func(m *mocks.MockUseCase) {
				m.OnWithdraw = func(ctx context.Context, id uint, amount decimal.Decimal) error {
					return nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:     "invalid wallet ID",
			walletID: "invalid",
			reqBody: WithdrawRequest{
				Amount: decimal.NewFromFloat(50.0),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid request body",
			walletID:   "1",
			reqBody:    "invalid json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "insufficient balance",
			walletID: "1",
			reqBody: WithdrawRequest{
				Amount: decimal.NewFromFloat(1000.0),
			},
			setupMock: func(m *mocks.MockUseCase) {
				m.OnWithdraw = func(ctx context.Context, id uint, amount decimal.Decimal) error {
					return errors.InsufficientBalance
				}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "wallet not found",
			walletID: "999",
			reqBody: WithdrawRequest{
				Amount: decimal.NewFromFloat(50.0),
			},
			setupMock: func(m *mocks.MockUseCase) {
				m.OnWithdraw = func(ctx context.Context, id uint, amount decimal.Decimal) error {
					return errors.RecordNotFound
				}
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "internal server error",
			walletID: "1",
			reqBody: WithdrawRequest{
				Amount: decimal.NewFromFloat(50.0),
			},
			setupMock: func(m *mocks.MockUseCase) {
				m.OnWithdraw = func(ctx context.Context, id uint, amount decimal.Decimal) error {
					return errors.InternalServer
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:     "negative amount",
			walletID: "1",
			reqBody: WithdrawRequest{
				Amount: decimal.NewFromFloat(-50.0),
			},
			setupMock: func(m *mocks.MockUseCase) {
				m.OnWithdraw = func(ctx context.Context, id uint, amount decimal.Decimal) error {
					return errors.InvalidArgs
				}
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mocks.MockUseCase{}
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			h := NewHandler(mockUC)
			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/wallets/"+tt.walletID+"/withdraw", bytes.NewReader(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.walletID})
			w := httptest.NewRecorder()

			h.Withdraw(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Withdraw() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_Transfer(t *testing.T) {
	tests := []struct {
		name       string
		walletID   string
		reqBody    interface{}
		setupMock  func(*mocks.MockUseCase)
		wantStatus int
	}{
		{
			name:     "successful transfer",
			walletID: "1",
			reqBody: TransferRequest{
				TargetWalletID: 2,
				Amount:         decimal.NewFromFloat(50.0),
			},
			setupMock: func(m *mocks.MockUseCase) {
				m.OnTransfer = func(ctx context.Context, fromID, toID uint, amount decimal.Decimal) error {
					return nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:     "invalid source wallet ID",
			walletID: "invalid",
			reqBody: TransferRequest{
				TargetWalletID: 2,
				Amount:         decimal.NewFromFloat(50.0),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid request body",
			walletID:   "1",
			reqBody:    "invalid json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "insufficient balance",
			walletID: "1",
			reqBody: TransferRequest{
				TargetWalletID: 2,
				Amount:         decimal.NewFromFloat(1000.0),
			},
			setupMock: func(m *mocks.MockUseCase) {
				m.OnTransfer = func(ctx context.Context, fromID, toID uint, amount decimal.Decimal) error {
					return errors.InsufficientBalance
				}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "source wallet not found",
			walletID: "999",
			reqBody: TransferRequest{
				TargetWalletID: 2,
				Amount:         decimal.NewFromFloat(50.0),
			},
			setupMock: func(m *mocks.MockUseCase) {
				m.OnTransfer = func(ctx context.Context, fromID, toID uint, amount decimal.Decimal) error {
					return errors.RecordNotFound
				}
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "target wallet not found",
			walletID: "1",
			reqBody: TransferRequest{
				TargetWalletID: 999,
				Amount:         decimal.NewFromFloat(50.0),
			},
			setupMock: func(m *mocks.MockUseCase) {
				m.OnTransfer = func(ctx context.Context, fromID, toID uint, amount decimal.Decimal) error {
					return errors.RecordNotFound
				}
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "internal server error",
			walletID: "1",
			reqBody: TransferRequest{
				TargetWalletID: 2,
				Amount:         decimal.NewFromFloat(50.0),
			},
			setupMock: func(m *mocks.MockUseCase) {
				m.OnTransfer = func(ctx context.Context, fromID, toID uint, amount decimal.Decimal) error {
					return errors.InternalServer
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:     "negative amount",
			walletID: "1",
			reqBody: TransferRequest{
				TargetWalletID: 2,
				Amount:         decimal.NewFromFloat(-50.0),
			},
			setupMock: func(m *mocks.MockUseCase) {
				m.OnTransfer = func(ctx context.Context, fromID, toID uint, amount decimal.Decimal) error {
					return errors.InvalidArgs
				}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "same wallet transfer",
			walletID: "1",
			reqBody: TransferRequest{
				TargetWalletID: 1,
				Amount:         decimal.NewFromFloat(50.0),
			},
			setupMock: func(m *mocks.MockUseCase) {
				m.OnTransfer = func(ctx context.Context, fromID, toID uint, amount decimal.Decimal) error {
					return errors.InvalidArgs
				}
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mocks.MockUseCase{}
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			h := NewHandler(mockUC)
			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/wallets/"+tt.walletID+"/transfer", bytes.NewReader(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.walletID})
			w := httptest.NewRecorder()

			h.Transfer(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Transfer() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_Balance(t *testing.T) {
	tests := []struct {
		name       string
		walletID   string
		setupMock  func(*mocks.MockUseCase)
		wantStatus int
		wantBody   string
	}{
		{
			name:     "successful balance retrieval",
			walletID: "1",
			setupMock: func(m *mocks.MockUseCase) {
				m.OnWallet = func(ctx context.Context, id uint) (*wallet.Wallet, error) {
					return &wallet.Wallet{
						ID:      1,
						Balance: decimal.NewFromFloat(100.50),
					}, nil
				}
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"balance":"100.5"}`,
		},
		{
			name:       "invalid wallet ID",
			walletID:   "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "wallet not found",
			walletID: "999",
			setupMock: func(m *mocks.MockUseCase) {
				m.OnWallet = func(ctx context.Context, id uint) (*wallet.Wallet, error) {
					return nil, errors.RecordNotFound
				}
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "internal server error",
			walletID: "1",
			setupMock: func(m *mocks.MockUseCase) {
				m.OnWallet = func(ctx context.Context, id uint) (*wallet.Wallet, error) {
					return nil, errors.InternalServer
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mocks.MockUseCase{}
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			h := NewHandler(mockUC)
			req := httptest.NewRequest(http.MethodGet, "/wallets/"+tt.walletID+"/balance", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.walletID})
			w := httptest.NewRecorder()

			h.Balance(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Balance() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantBody != "" {
				if body := w.Body.String(); body != tt.wantBody+"\n" {
					t.Errorf("Balance() body = %v, want %v", body, tt.wantBody)
				}
			}
		})
	}
}

func TestHandler_Transactions(t *testing.T) {
	mockTxs := []transaction.Transaction{
		{
			ID:         1,
			Method:     transaction.MethodDeposit,
			TxAt:       time.Now(),
			Amount:     decimal.NewFromFloat(100),
			ToWalletID: 1,
		},
		{
			ID:           2,
			Method:       transaction.MethodWithdraw,
			TxAt:         time.Now(),
			Amount:       decimal.NewFromFloat(50),
			FromWalletID: 1,
		},
	}

	tests := []struct {
		name       string
		walletID   string
		setupMock  func(*mocks.MockUseCase)
		wantStatus int
		wantTxs    []transaction.Transaction
	}{
		{
			name:     "successful transactions retrieval",
			walletID: "1",
			setupMock: func(m *mocks.MockUseCase) {
				m.OnWalletTransactions = func(ctx context.Context, id uint) ([]transaction.Transaction, error) {
					return mockTxs, nil
				}
			},
			wantStatus: http.StatusOK,
			wantTxs:    mockTxs,
		},
		{
			name:       "invalid wallet ID",
			walletID:   "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "wallet not found",
			walletID: "999",
			setupMock: func(m *mocks.MockUseCase) {
				m.OnWalletTransactions = func(ctx context.Context, id uint) ([]transaction.Transaction, error) {
					return nil, errors.RecordNotFound
				}
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "internal server error",
			walletID: "1",
			setupMock: func(m *mocks.MockUseCase) {
				m.OnWalletTransactions = func(ctx context.Context, id uint) ([]transaction.Transaction, error) {
					return nil, errors.InternalServer
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:     "empty transaction list",
			walletID: "1",
			setupMock: func(m *mocks.MockUseCase) {
				m.OnWalletTransactions = func(ctx context.Context, id uint) ([]transaction.Transaction, error) {
					return []transaction.Transaction{}, nil
				}
			},
			wantStatus: http.StatusOK,
			wantTxs:    []transaction.Transaction{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mocks.MockUseCase{}
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			h := NewHandler(mockUC)
			req := httptest.NewRequest(http.MethodGet, "/wallets/"+tt.walletID+"/transactions", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.walletID})
			w := httptest.NewRecorder()

			h.Transactions(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Transactions() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantTxs != nil {
				var gotTxs []transaction.Transaction
				if err := json.NewDecoder(w.Body).Decode(&gotTxs); err != nil {
					t.Errorf("Failed to decode response body: %v", err)
				}
				if len(gotTxs) != len(tt.wantTxs) {
					t.Errorf("Transactions() returned %d transactions, want %d", len(gotTxs), len(tt.wantTxs))
				}
			}
		})
	}
}
