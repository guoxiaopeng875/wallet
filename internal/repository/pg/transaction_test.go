package pg

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/wallet/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransactionRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defaultConnTestRunner.RunTest(ctx, t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		tp := NewTransactionRepository(NewRepository(conn))
		tx := &transaction.Transaction{
			Method:       transaction.MethodTransfer,
			TxAt:         time.Date(2024, 11, 5, 0, 0, 0, 0, time.Local),
			Amount:       decimal.NewFromFloat(100.1111),
			FromWalletID: 1,
			ToWalletID:   10,
		}
		err := tp.Create(ctx, tx)
		assert.NoError(t, err)
		list, err := tp.ListByWalletID(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
		tx.ID = list[0].ID
		assert.Equal(t, *tx, list[0])
	})
}
