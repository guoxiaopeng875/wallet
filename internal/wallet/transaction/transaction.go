package transaction

import (
	"github.com/shopspring/decimal"
	"time"
)

// Method of transaction
type Method string

const (
	MethodDeposit  Method = "deposit"
	MethodWithdraw Method = "withdraw"
	MethodTransfer Method = "transfer"
)

type Transaction struct {
	ID           uint            `json:"id"`
	Method       Method          `json:"method"`
	TxAt         time.Time       `json:"tx_at"`
	Amount       decimal.Decimal `json:"amount"`
	FromWalletID uint            `json:"from_wallet_id"`
	ToWalletID   uint            `json:"to_wallet_id"`
}
