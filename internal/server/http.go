package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/guoxiaopeng875/wallet/internal/config"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// httpServer represents the HTTP server for the wallet API
type httpServer struct {
	*http.Server
	listener net.Listener
	addr     string
}

// NewServer creates a new HTTP server instance
func NewServer(h *Handler, conf *config.Config) Server {
	router := mux.NewRouter()
	router.Use(LoggingMiddleware())

	// Register routes
	router.HandleFunc("/wallets/{id}/deposit", h.Deposit).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/withdraw", h.Withdraw).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/transfer", h.Transfer).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/balance", h.Balance).Methods(http.MethodGet)
	router.HandleFunc("/wallets/{id}/transactions", h.Transactions).Methods(http.MethodGet)

	// Add health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	srv := &httpServer{
		Server: &http.Server{
			Handler: router,
			Addr:    conf.Server.Address,
		},
	}

	return srv
}

// Start starts the HTTP server
func (s *httpServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	s.listener = listener

	logrus.Infof("HTTP server listening on %s", s.Addr)
	if err := s.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

// Stop gracefully stops the HTTP server
func (s *httpServer) Stop(ctx context.Context) error {
	logrus.Info("Stopping HTTP server")
	return s.Shutdown(ctx)
}
