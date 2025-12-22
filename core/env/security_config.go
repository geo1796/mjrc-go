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
}

func loadSecurityConfig(isProd bool) (SecurityConfig, error) {
	cfg := SecurityConfig{JwtCookieName: getEnv("JWT_COOKIE_NAME", "__Host-jwt_")}

	if jwtSecret := getEnv("JWT_SECRET", ""); isProd && jwtSecret == "" {
		return SecurityConfig{}, errors.New("JWT_SECRET is not set")
	} else {
		cfg.JwtSecret = []byte(jwtSecret)
	}

	if jwtTTL, err := time.ParseDuration(getEnv("JWT_TTL", "1h")); err != nil {
		return SecurityConfig{}, fmt.Errorf("failed to parse JWT_TTL: %v", err)
	} else {
		cfg.JwtTTL = jwtTTL
	}

	return cfg, nil
}
