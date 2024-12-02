package config

import (
	"path"

	"github.com/ozgur-yalcin/mfa/src/backend"
)

type Config struct {
	DatabaseBackend backend.Backend
}

const (
	sqliteFileName = "mfa.db"
)

func DefaultConfig() *Config {
	return &Config{DatabaseBackend: backend.NewSqlite(path.Join(".", sqliteFileName))}
}