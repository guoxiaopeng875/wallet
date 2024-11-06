package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/guoxiaopeng875/wallet/internal/config"
	"github.com/guoxiaopeng875/wallet/internal/repository/pg"
	"github.com/guoxiaopeng875/wallet/internal/server"
	"github.com/guoxiaopeng875/wallet/internal/wallet"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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

	// Setup application
	app, cleanup, err := setupApp(conf)
	if err != nil {
		logrus.Fatalf("Failed to setup application: %v", err)
	}
	defer cleanup()

	// Run application
	if err := run(app); err != nil {
		logrus.Fatalf("Application error: %v", err)
	}
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

func setupApp(conf *config.Config) (server.Server, func(), error) {
	ctx := context.Background()

	// Initialize database
	conn, dbCloser, err := pg.NewConnect(ctx, conf.Repository.DSN)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize repositories and use cases
	repo := pg.NewRepository(conn)
	uc := wallet.NewUseCase(
		pg.NewWalletRepository(repo),
		pg.NewTransactionRepository(repo),
		pg.NewDBTx(repo),
	)

	// Initialize server
	srv := server.NewServer(
		server.NewHandler(uc),
		conf,
	)

	cleanup := func() {
		dbCloser()
	}

	return srv, cleanup, nil
}

func run(srv server.Server) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Error channel for server errors
	errCh := make(chan error, 1)
	go func() {
		if err := srv.Start(ctx); err != nil {
			errCh <- err
		}
	}()

	// Signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	// Wait for signal or error
	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case sig := <-sigCh:
		logrus.Infof("Received signal: %v", sig)
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	return nil
}
