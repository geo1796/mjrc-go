package env

import (
	"errors"
	"fmt"
	"time"
)

type SecurityConfig struct {
	AccessCookieName  string
	AccessTokenSecret []byte
	AccessTokenTTL    time.Duration
	AdminPassword     string
	APIKey            string
}

func loadSecurityConfig(isProd bool) (SecurityConfig, error) {
	cfg := SecurityConfig{
		AccessCookieName: getEnv("ACCESS_COOKIE_NAME", "__Host-jwt_"),
	}

	if jwtTTL, err := time.ParseDuration(getEnv("ACCESS_TOKEN_TTL", "1h")); err != nil {
		return SecurityConfig{}, fmt.Errorf("failed to parse ACCESS_TOKEN_TTL: %w", err)
	} else {
		cfg.AccessTokenTTL = jwtTTL
	}

	if isProd {
		if jwtSecret := getEnv("ACCESS_TOKEN_SECRET", ""); jwtSecret == "" {
			return SecurityConfig{}, errors.New("ACCESS_TOKEN_SECRET is not set")
		} else {
			cfg.AccessTokenSecret = []byte(jwtSecret)
		}

		if cfg.AdminPassword = getEnv("ADMIN_PASSWORD", ""); cfg.AdminPassword == "" {
			return SecurityConfig{}, errors.New("ADMIN_PASSWORD is not set")
		}

		if cfg.APIKey = getEnv("API_KEY", ""); cfg.APIKey == "" {
			return SecurityConfig{}, errors.New("API_KEY is not set")
		}
	}

	return cfg, nil
}
