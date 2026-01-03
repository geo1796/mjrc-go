package env

import (
	"errors"
	"fmt"
	"time"
)

type SecurityConfig struct {
	JwtCookieName string
	JwtSecret     []byte
	JwtTTL        time.Duration
	AdminPassword string
	APIKey        string
}

func loadSecurityConfig(isProd bool) (SecurityConfig, error) {
	cfg := SecurityConfig{
		JwtCookieName: getEnv("JWT_COOKIE_NAME", "__Host-jwt_"),
	}

	if jwtTTL, err := time.ParseDuration(getEnv("JWT_TTL", "1h")); err != nil {
		return SecurityConfig{}, fmt.Errorf("failed to parse JWT_TTL: %v", err)
	} else {
		cfg.JwtTTL = jwtTTL
	}

	if isProd {
		if jwtSecret := getEnv("JWT_SECRET", ""); jwtSecret == "" {
			return SecurityConfig{}, errors.New("JWT_SECRET is not set")
		} else {
			cfg.JwtSecret = []byte(jwtSecret)
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
