package env

import (
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

func loadPostgresConfig() (PostgresConfig, error) {
	cfg := PostgresConfig{DSN: getEnv("PG_DSN", "")}

	var err error

	if cfg.ConnMaxLifetime, err = time.ParseDuration(getEnv("PG_CONN_MAX_LIFETIME", "10m")); err != nil {
		return PostgresConfig{}, fmt.Errorf("failed to parse PG_CONN_MAX_LIFETIME: %v", err)
	}

	if cfg.ConnMaxIdleTime, err = time.ParseDuration(getEnv("PG_CONN_MAX_IDLE_TIME", "5m")); err != nil {
		return PostgresConfig{}, fmt.Errorf("failed to parse PG_CONN_MAX_IDLE_TIME: %v", err)
	}

	if cfg.MaxOpenConns, err = strconv.Atoi(getEnv("PG_MAX_OPEN_CONNS", "2")); err != nil {
		return PostgresConfig{}, fmt.Errorf("failed to parse PG_MAX_OPEN_CONNS: %v", err)
	}

	if cfg.MaxIdleConns, err = strconv.Atoi(getEnv("PG_MAX_IDLE_CONNS", "1")); err != nil {
		return PostgresConfig{}, fmt.Errorf("failed to parse PG_MAX_IDLE_CONNS: %v", err)
	}

	return cfg, nil
}
