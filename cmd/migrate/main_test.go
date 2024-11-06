package main

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/config"
	"github.com/guoxiaopeng875/wallet/internal/repository/pg"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var defaultTestDB string

func init() {
	defaultTestDB = os.Getenv("PGX_TEST_DATABASE")
	if defaultTestDB == "" {
		defaultTestDB = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}
}

func TestSetupLogger(t *testing.T) {
	setupLogger()
	assert.IsType(t, &logrus.JSONFormatter{}, logrus.StandardLogger().Formatter)
	assert.Equal(t, logrus.InfoLevel, logrus.StandardLogger().Level)
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "valid path",
			path:    "testdata/test.json",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadConfig(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRunMigration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conf := &config.Config{
		Repository: config.Repository{
			MigrateDSN: defaultTestDB,
		},
	}

	err := runMigration(conf)
	require.NoError(t, err)

	// 验证数据库和表是否创建成功
	conn, closer, err := pg.NewConnect(ctx, defaultTestDB)
	require.NoError(t, err)
	defer closer()

	var exists bool
	// 检查表是否存在
	tables := []string{"wallets", "transactions"}
	for _, table := range tables {
		err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = $1)", table).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "%s table should exist", table)
	}
}

func TestRunMigrationWithInvalidDSN(t *testing.T) {
	conf := &config.Config{
		Repository: config.Repository{
			MigrateDSN: "invalid://dsn",
		},
	}

	err := runMigration(conf)
	assert.Error(t, err)
}
