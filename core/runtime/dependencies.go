package runtime

import (
	"mjrc/core/env"
	"mjrc/core/postgres"
	"mjrc/core/security"
)

type Dependencies interface {
	APIConfig() env.APIConfig
	DB() postgres.DB
	JWT() security.JWT
}

func New(apiConfig env.APIConfig, db postgres.DB, jwt security.JWT) Dependencies {
	return &dependencies{
		apiConfig: apiConfig,
		db:        db,
		jwt:       jwt,
	}
}

type dependencies struct {
	apiConfig env.APIConfig
	db        postgres.DB
	jwt       security.JWT
}

func (d *dependencies) APIConfig() env.APIConfig {
	return d.apiConfig
}

func (d *dependencies) DB() postgres.DB {
	return d.db
}

func (d *dependencies) JWT() security.JWT {
	return d.jwt
}
