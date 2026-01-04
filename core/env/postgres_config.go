package env

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type PostgresConfig struct {
	DSN             string
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
}

func loadPostgresConfig(isProd bool) (PostgresConfig, error) {
	cfg := PostgresConfig{DSN: getEnv("PG_DSN", "")}

	if isProd && cfg.DSN == "" {
		return PostgresConfig{}, errors.New("PG_DSN is not set")
	}

	var err error

	if cfg.ConnMaxLifetime, err = time.ParseDuration(getEnv("PG_CONN_MAX_LIFETIME", "10m")); err != nil {
		return PostgresConfig{}, fmt.Errorf("failed to parse PG_CONN_MAX_LIFETIME: %w", err)
	}

	if cfg.ConnMaxIdleTime, err = time.ParseDuration(getEnv("PG_CONN_MAX_IDLE_TIME", "5m")); err != nil {
		return PostgresConfig{}, fmt.Errorf("failed to parse PG_CONN_MAX_IDLE_TIME: %w", err)
	}

	if cfg.MaxOpenConns, err = strconv.Atoi(getEnv("PG_MAX_OPEN_CONNS", "2")); err != nil {
		return PostgresConfig{}, fmt.Errorf("failed to parse PG_MAX_OPEN_CONNS: %w", err)
	}

	if cfg.MaxIdleConns, err = strconv.Atoi(getEnv("PG_MAX_IDLE_CONNS", "1")); err != nil {
		return PostgresConfig{}, fmt.Errorf("failed to parse PG_MAX_IDLE_CONNS: %w", err)
	}

	return cfg, nil
}
