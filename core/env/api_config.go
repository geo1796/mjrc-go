package env

import "fmt"

type APIConfig struct {
	Address       string
	ApiKey        string
	AdminPassword string
}

func loadAPIConfig(isProd bool) (APIConfig, error) {
	apiKey := getEnv("API_KEY", "")
	if apiKey == "" && isProd {
		return APIConfig{}, fmt.Errorf("API_KEY is not set")
	}

	adminPassword := getEnv("ADMIN_PASSWORD", "")
	if adminPassword == "" && isProd {
		return APIConfig{}, fmt.Errorf("ADMIN_PASSWORD is not set")
	}

	return APIConfig{
		Address:       getEnv("API_ADDRESS", ":8080"),
		ApiKey:        apiKey,
		AdminPassword: adminPassword,
	}, nil
}
