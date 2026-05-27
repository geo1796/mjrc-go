package env

import (
	"errors"
	"fmt"
	"time"
)

type ServerConfig struct {
	Addr              string
	AccessTokenSecret []byte
	AccessTokenTTL    time.Duration
	AdminPassword     string
	APIKey            string
}

func loadServerConfig(isProd bool) (ServerConfig, error) {
	cfg := ServerConfig{Addr: "0.0.0.0:" + getEnv("PORT", "8080")}

	if accessTokenTTL, err := time.ParseDuration(getEnv("ACCESS_TOKEN_TTL", "1h")); err != nil {
		return ServerConfig{}, fmt.Errorf("failed to parse ACCESS_TOKEN_TTL: %w", err)
	} else {
		cfg.AccessTokenTTL = accessTokenTTL
	}

	if accessTokenSecret := getEnv("ACCESS_TOKEN_SECRET", ""); isProd && accessTokenSecret == "" {
		return ServerConfig{}, errors.New("ACCESS_TOKEN_SECRET is not set")
	} else {
		cfg.AccessTokenSecret = []byte(accessTokenSecret)
	}

	if cfg.AdminPassword = getEnv("ADMIN_PASSWORD", ""); isProd && cfg.AdminPassword == "" {
		return ServerConfig{}, errors.New("ADMIN_PASSWORD is not set")
	}

	if cfg.APIKey = getEnv("API_KEY", ""); isProd && cfg.APIKey == "" {
		return ServerConfig{}, errors.New("API_KEY is not set")
	}

	return cfg, nil
}
