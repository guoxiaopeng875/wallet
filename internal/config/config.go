package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Repository Repository `json:"repository"`
	Server     Server     `json:"server"`
}

type Repository struct {
	DSN        string `json:"dsn"`
	MigrateDSN string `json:"migrate_dsn"`
}

type Server struct {
	Address string `json:"address"`
}

func NewConfig(confFile string) (*Config, error) {
	f, err := os.Open(confFile)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := json.NewDecoder(f).Decode(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
