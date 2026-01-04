package env

import (
	"errors"
	"fmt"
	"time"
)

type SecurityConfig struct {
	AccessTokenSecret []byte
	AccessTokenTTL    time.Duration
	AdminPassword     string
	APIKey            string
}

func loadSecurityConfig(isProd bool) (SecurityConfig, error) {
	var cfg SecurityConfig

	if accessTokenTTL, err := time.ParseDuration(getEnv("ACCESS_TOKEN_TTL", "1h")); err != nil {
		return SecurityConfig{}, fmt.Errorf("failed to parse ACCESS_TOKEN_TTL: %w", err)
	} else {
		cfg.AccessTokenTTL = accessTokenTTL
	}

	if accessTokenSecret := getEnv("ACCESS_TOKEN_SECRET", ""); isProd && accessTokenSecret == "" {
		return SecurityConfig{}, errors.New("ACCESS_TOKEN_SECRET is not set")
	} else {
		cfg.AccessTokenSecret = []byte(accessTokenSecret)
	}

	if cfg.AdminPassword = getEnv("ADMIN_PASSWORD", ""); isProd && cfg.AdminPassword == "" {
		return SecurityConfig{}, errors.New("ADMIN_PASSWORD is not set")
	}

	if cfg.APIKey = getEnv("API_KEY", ""); isProd && cfg.APIKey == "" {
		return SecurityConfig{}, errors.New("API_KEY is not set")
	}

	return cfg, nil
}
