package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantErr  bool
		validate func(*testing.T, *Config)
	}{
		{
			name: "valid config",
			content: `{
				"repository": {
					"dsn": "postgres://user:pass@localhost:5432/db",
					"migrate_dsn": "postgres://user:pass@localhost:5432/db?sslmode=disable"
				},
				"server": {
					"address": ":8080"
				}
			}`,
			wantErr: false,
			validate: func(t *testing.T, c *Config) {
				if c.Repository.DSN != "postgres://user:pass@localhost:5432/db" {
					t.Errorf("expected DSN %s, got %s", "postgres://user:pass@localhost:5432/db", c.Repository.DSN)
				}
				if c.Repository.MigrateDSN != "postgres://user:pass@localhost:5432/db?sslmode=disable" {
					t.Errorf("expected MigrateDSN %s, got %s", "postgres://user:pass@localhost:5432/db?sslmode=disable", c.Repository.MigrateDSN)
				}
				if c.Server.Address != ":8080" {
					t.Errorf("expected Address %s, got %s", ":8080", c.Server.Address)
				}
			},
		},
		{
			name: "invalid json",
			content: `{
				"repository": {
					"dsn": "postgres://user:pass@localhost:5432/db",
				}
			}`,
			wantErr: true,
		},
		{
			name:     "empty config",
			content:  "{}",
			wantErr:  false,
			validate: func(t *testing.T, c *Config) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "config.json")

			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("failed to write test config file: %v", err)
			}

			// Test config loading
			cfg, err := NewConfig(tmpFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestNewConfig_FileError(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "non-existent file",
			path:    "/non/existent/path/config.json",
			wantErr: true,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfig(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
