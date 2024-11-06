package main

import (
	"context"
	"github.com/guoxiaopeng875/wallet/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestSetupLogger(t *testing.T) {
	setupLogger()
	// Verify logger settings
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
			name:    "invalid path",
			path:    "nonexistent.json",
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

type mockServer struct {
	startErr error
	stopErr  error
}

func (m *mockServer) Start(ctx context.Context) error {
	return m.startErr
}

func (m *mockServer) Stop(ctx context.Context) error {
	return m.stopErr
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		server  *mockServer
		signal  os.Signal
		wantErr bool
	}{
		{
			name:    "successful run and shutdown",
			server:  &mockServer{},
			signal:  syscall.SIGTERM,
			wantErr: false,
		},
		{
			name: "server start error",
			server: &mockServer{
				startErr: assert.AnError,
			},
			wantErr: true,
		},
		{
			name: "server stop error",
			server: &mockServer{
				stopErr: assert.AnError,
			},
			signal:  syscall.SIGTERM,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errCh := make(chan error, 1)

			go func() {
				errCh <- run(tt.server)
			}()

			if tt.signal != nil {
				time.Sleep(100 * time.Millisecond)
				syscall.Kill(syscall.Getpid(), tt.signal.(syscall.Signal))
			}

			err := <-errCh
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetupApp(t *testing.T) {
	tests := []struct {
		name    string
		conf    *config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			conf: &config.Config{
				Repository: config.Repository{
					DSN: os.Getenv("PGX_TEST_DATABASE"),
				},
				Server: config.Server{
					Address: ":8080",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, cleanup, err := setupApp(tt.conf)
			defer cleanup()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, app)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, app)
			}
		})
	}
}
