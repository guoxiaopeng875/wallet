package pg

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/wallet"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWalletRepository_Get(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defaultConnTestRunner.RunTest(ctx, t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		id := uint(1)
		wp := NewWalletRepository(NewRepository(conn))
		w, err := wp.Get(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, w)

		mustExec(ctx, t, conn, "insert into wallets (balance) values (100.1122);")
		w = mustGetWallet(ctx, t, wp, id)
		assert.Equal(t, w, &wallet.Wallet{
			ID:      id,
			Balance: decimal.NewFromFloat(100.1122),
		})
	})
}

func TestWalletRepository_UpdateBalance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defaultConnTestRunner.RunTest(ctx, t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		id := uint(1)
		wp := NewWalletRepository(NewRepository(conn))
		mustExec(ctx, t, conn, "insert into wallets (balance) values (0.0000);")
		w := mustGetWallet(ctx, t, wp, id)
		// add balance
		err := wp.UpdateBalance(ctx, w, decimal.NewFromFloat(1))
		assert.NoError(t, err)

		w = mustGetWallet(ctx, t, wp, id)
		assert.Equal(t, "1", w.Balance.String())

		// add balance
		err = wp.UpdateBalance(ctx, w, decimal.NewFromFloat(2.2222))
		assert.NoError(t, err)

		w = mustGetWallet(ctx, t, wp, id)
		assert.Equal(t, "3.2222", w.Balance.String())

		// sub balance
		err = wp.UpdateBalance(ctx, w, decimal.NewFromFloat(-3.2))
		assert.NoError(t, err)

		w = mustGetWallet(ctx, t, wp, id)
		assert.Equal(t, "0.0222", w.Balance.String())

		// update failed
		w.ID = 0
		err = wp.UpdateBalance(ctx, w, decimal.NewFromFloat(-3.2))
		assert.Error(t, err)

	})
}

func TestWalletRepository_UpdateConcurrently(t *testing.T) {
	t.Skip()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defaultConnTestRunner.RunTest(ctx, t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		id := uint(1)
		wp := NewWalletRepository(NewRepository(conn))
		mustExec(ctx, t, conn, "insert into wallets (balance) values (0.0000);")
		w := mustGetWallet(ctx, t, wp, id)
		var wg sync.WaitGroup
		errCount := &atomic.Uint32{}
		concurrentCount := 50
		for i := 0; i < concurrentCount; i++ {
			go func() {
				wg.Add(1)
				err := wp.UpdateBalance(ctx, w, decimal.NewFromFloat(1))
				if err == nil {
					w = mustGetWallet(ctx, t, wp, id)
					assert.Equal(t, "1", w.Balance.String())
				} else {
					errCount.Add(1)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		// there can only be one success
		assert.Equal(t, uint32(concurrentCount-1), errCount.Load())
	})
}

func BenchmarkWalletRepository_Update(b *testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defaultConnTestRunner.RunTest(ctx, b, func(ctx context.Context, _ testing.TB, conn *pgx.Conn) {
		b.ResetTimer()

		wp := NewWalletRepository(NewRepository(conn))
		for i := 0; i < b.N; i++ {
			id := uint(i + 1)
			mustExec(ctx, b, conn, "insert into wallets (id, balance) values ($1, 0.0000);", id)
			w := mustGetWallet(ctx, b, wp, id)
			// add balance
			err := wp.UpdateBalance(ctx, w, decimal.NewFromFloat(2.2222))
			assert.NoError(b, err)
			w1 := mustGetWallet(ctx, b, wp, id)
			assert.Equal(b, "2.2222", w1.Balance.String())
		}
	})
}

func mustGetWallet(ctx context.Context, t testing.TB, wp wallet.Repository, id uint) *wallet.Wallet {
	w, err := wp.Get(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, w)
	return w
}
