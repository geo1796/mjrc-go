package env

import (
	"fmt"
	"os"
)

func getEnv(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}

type Env interface {
	IsProd() bool
	ActiveProfile() string
	PostgresConfig() PostgresConfig
	ServerConfig() ServerConfig
}

type env struct {
	activeProfile  string
	postgresConfig PostgresConfig
	serverConfig   ServerConfig
}

func (e *env) IsProd() bool {
	return e.activeProfile == "prod"
}

func (e *env) ActiveProfile() string {
	return e.activeProfile
}

func (e *env) PostgresConfig() PostgresConfig {
	return e.postgresConfig
}

func (e *env) ServerConfig() ServerConfig {
	return e.serverConfig
}

func Load() (Env, error) {
	e := &env{
		activeProfile: getEnv("ACTIVE_PROFILE", "test"),
	}

	var err error

	if e.postgresConfig, err = loadPostgresConfig(e.IsProd()); err != nil {
		return nil, fmt.Errorf("failed to load postgres config: %w", err)
	}

	if e.serverConfig, err = loadServerConfig(e.IsProd()); err != nil {
		return nil, fmt.Errorf("failed to load security config: %w", err)
	}

	return e, nil
}
