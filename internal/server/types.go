package server

import "github.com/shopspring/decimal"

// Request types for API endpoints
type (
	DepositRequest struct {
		Amount decimal.Decimal `json:"amount" validate:"required,gt=0"`
	}

	WithdrawRequest struct {
		Amount decimal.Decimal `json:"amount" validate:"required,gt=0"`
	}

	TransferRequest struct {
		TargetWalletID uint            `json:"target_wallet_id" validate:"required,gt=0"`
		Amount         decimal.Decimal `json:"amount" validate:"required,gt=0"`
	}

	// Response types
	BalanceResponse struct {
		Balance string `json:"balance"`
	}
)
