package env

import (
	"errors"
)

type APIConfig struct {
	Address string
	ApiKey  string
}

func loadAPIConfig(isProd bool) (APIConfig, error) {
	apiKey := getEnv("API_KEY", "")
	if apiKey == "" && isProd {
		return APIConfig{}, errors.New("API_KEY is not set")
	}

	return APIConfig{
		Address: getEnv("API_ADDRESS", ":8080"),
		ApiKey:  apiKey,
	}, nil
}
