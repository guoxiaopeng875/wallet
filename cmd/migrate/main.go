package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/guoxiaopeng875/wallet/internal/config"
	"github.com/guoxiaopeng875/wallet/internal/repository/pg"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
)

type migrator struct {
	ctx    context.Context
	conn   *pgx.Conn
	closer func()
}

func main() {
	// Parse command line flags
	configPath := flag.String("conf", "", "config path, eg: -conf config.json")
	flag.Parse()

	// Initialize logger
	setupLogger()

	// Load configuration
	conf, err := loadConfig(*configPath)
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	// Run migration
	if err := runMigration(conf); err != nil {
		logrus.Fatalf("Migration failed: %v", err)
	}

	logrus.Info("Migration completed successfully")
}

func setupLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

func loadConfig(path string) (*config.Config, error) {
	if path == "" {
		return nil, fmt.Errorf("config path is required")
	}
	return config.NewConfig(path)
}

func runMigration(conf *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Connect to database
	conn, closer, err := pg.NewConnect(ctx, conf.Repository.MigrateDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer closer()

	m := &migrator{
		ctx:    ctx,
		conn:   conn,
		closer: closer,
	}

	// Run migrations in order
	migrations := []struct {
		name string
		fn   func() error
	}{
		{"Create wallet table", m.createWalletTable},
		{"Create transaction table", m.createTransactionTable},
		{"Insert initial data", m.insertInitialData},
	}

	for _, migration := range migrations {
		logrus.Infof("Running migration: %s", migration.name)
		if err := migration.fn(); err != nil {
			return fmt.Errorf("migration '%s' failed: %w", migration.name, err)
		}
		logrus.Infof("Completed migration: %s", migration.name)
	}

	return nil
}

func (m *migrator) createWalletTable() error {
	tx, err := m.conn.Begin(m.ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(m.ctx)

	query := `
		CREATE TABLE IF NOT EXISTS wallets (
			id SERIAL PRIMARY KEY,
			balance DECIMAL(20,4) NOT NULL DEFAULT 0.0000
		);
		ALTER TABLE IF EXISTS public.wallets OWNER to postgres;
	`

	if _, err := tx.Exec(m.ctx, query); err != nil {
		return fmt.Errorf("failed to create wallet table: %w", err)
	}

	return tx.Commit(m.ctx)
}

func (m *migrator) createTransactionTable() error {
	tx, err := m.conn.Begin(m.ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(m.ctx)

	query := `
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			method VARCHAR(10) NOT NULL,
			tx_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			amount DECIMAL(20,4) NOT NULL,
			from_wallet_id INTEGER,
			to_wallet_id INTEGER
		);
		ALTER TABLE IF EXISTS public.transactions OWNER to postgres;
	`

	if _, err := tx.Exec(m.ctx, query); err != nil {
		return fmt.Errorf("failed to create transaction table: %w", err)
	}

	return tx.Commit(m.ctx)
}

func (m *migrator) insertInitialData() error {
	tx, err := m.conn.Begin(m.ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(m.ctx)

	query := `
		INSERT INTO wallets (balance) VALUES
		(0.00),
		(0.00),
		(0.00),
		(0.00),
		(0.00);
	`

	if _, err := tx.Exec(m.ctx, query); err != nil {
		return fmt.Errorf("failed to insert initial data: %w", err)
	}

	return tx.Commit(m.ctx)
}
