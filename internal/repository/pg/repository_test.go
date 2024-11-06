package pg

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxtest"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var defaultConnTestRunner pgxtest.ConnTestRunner

func init() {
	logrus.SetLevel(logrus.ErrorLevel)
	defaultConnTestRunner = pgxtest.DefaultConnTestRunner()
	defaultConnTestRunner.CreateConfig = func(ctx context.Context, t testing.TB) *pgx.ConnConfig {
		config, err := pgx.ParseConfig(os.Getenv("PGX_TEST_DATABASE"))
		require.NoError(t, err)
		return config
	}
	defaultConnTestRunner.AfterConnect = func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		mustExec(ctx, t, conn, `CREATE TEMPORARY TABLE wallets (
		id SERIAL PRIMARY KEY,
		balance DECIMAL(20,4) NOT NULL DEFAULT 0.0000
		)`)
		mustExec(ctx, t, conn, `CREATE TEMPORARY TABLE transactions (
		id SERIAL PRIMARY KEY,
		method VARCHAR(10) NOT NULL,
		tx_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		amount DECIMAL(20,4) NOT NULL,
		from_wallet_id INTEGER,
		to_wallet_id INTEGER
		)`)
	}
}

func mustExec(ctx context.Context, t testing.TB, conn *pgx.Conn, sql string, arguments ...any) (commandTag pgconn.CommandTag) {
	var err error
	if commandTag, err = conn.Exec(ctx, sql, arguments...); err != nil {
		t.Fatalf("Exec unexpectedly failed with %v: %v", sql, err)
	}
	return
}

func TestNewConnect(t *testing.T) {
	db, cleanup, err := NewConnect(context.Background(), os.Getenv("PGX_TEST_DATABASE"))
	assert.NoError(t, err)
	defer cleanup()
	assert.NotNil(t, db)
}

func TestNewDBTx(t *testing.T) {
	assert.Nil(t, NewDBTx(nil))
}

func TestExecTx(t *testing.T) {
	ctx := context.Background()
	defaultConnTestRunner.RunTest(ctx, t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		repo := NewRepository(conn)
		dbTx := NewDBTx(repo)
		assert.NotNil(t, dbTx)
		err := dbTx.ExecTx(ctx, func(ctx context.Context) error {
			mustExec(ctx, t, repo.DB(ctx), "insert into wallets (balance) values (100.1122);")
			mustExec(ctx, t, repo.DB(ctx), "insert into wallets (balance) values (100.1122);")
			return nil
		})
		assert.NoError(t, err)
		var count int
		err = conn.QueryRow(ctx, "select count(1) from wallets;").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})
}

func TestExecTxFailed(t *testing.T) {
	ctx := context.Background()
	defaultConnTestRunner.RunTest(ctx, t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		repo := NewRepository(conn)
		dbTx := NewDBTx(repo)
		assert.NotNil(t, dbTx)
		err := dbTx.ExecTx(ctx, func(ctx context.Context) error {
			mustExec(ctx, t, repo.DB(ctx), "insert into wallets (balance) values (100.1122);")
			return errors.New(1, "mock failed")
		})
		assert.Error(t, err)
		var count int
		err = conn.QueryRow(ctx, "select count(1) from wallets;").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}
