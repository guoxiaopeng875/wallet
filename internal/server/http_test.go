package server

import (
	"context"
	"fmt"
	"github.com/guoxiaopeng875/wallet/internal/config"
	"github.com/guoxiaopeng875/wallet/internal/server/mocks"
	"net"
	"net/http"
	"testing"
	"time"
)

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func TestServer_StartStop(t *testing.T) {
	// Get a free port
	port, err := getFreePort()
	if err != nil {
		t.Fatalf("Failed to get free port: %v", err)
	}

	// Setup mock and handler
	mockUC := &mocks.MockUseCase{}
	h := NewHandler(mockUC)

	// Create server with the free port
	conf := &config.Config{
		Server: config.Server{
			Address: fmt.Sprintf(":%d", port),
		},
	}
	srv := NewServer(h, conf)

	// Create context with timeout for the entire test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start(ctx)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Verify server is running
	client := &http.Client{
		Timeout: time.Second,
	}
	resp, err := client.Get(fmt.Sprintf("http://localhost:%d/health", port))
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	resp.Body.Close()

	// Stop server
	if err := srv.Stop(ctx); err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}

	// Verify server has stopped
	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server.Start() returned unexpected error: %v", err)
		}
	case <-ctx.Done():
		t.Error("Context timeout while waiting for server to stop")
	}
}
