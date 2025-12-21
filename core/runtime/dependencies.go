package runtime

import (
	"mjrc/core/env"
	"mjrc/core/postgres"
)

type Dependencies interface {
	APIConfig() env.APIConfig
	DB() postgres.DB
}

func New(apiConfig env.APIConfig, db postgres.DB) Dependencies {
	return &dependencies{apiConfig, db}
}

type dependencies struct {
	apiConfig env.APIConfig
	db        postgres.DB
}

func (d *dependencies) APIConfig() env.APIConfig {
	return d.apiConfig
}

func (d *dependencies) DB() postgres.DB {
	return d.db
}
