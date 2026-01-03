package runtime

import (
	"mjrc/core/postgres"
	"mjrc/core/security"
)

type Dependencies interface {
	DB() postgres.DB
	JWT() security.JWT
	AdminAuthenticator() security.Authenticator
	APIKeyAuthenticator() security.Authenticator
}

type dependencies struct {
	db                  postgres.DB
	jwt                 security.JWT
	adminAuthenticator  security.Authenticator
	apiKeyAuthenticator security.Authenticator
}

func (d *dependencies) DB() postgres.DB {
	return d.db
}

func (d *dependencies) JWT() security.JWT {
	return d.jwt
}

func (d *dependencies) AdminAuthenticator() security.Authenticator {
	return d.adminAuthenticator
}

func (d *dependencies) APIKeyAuthenticator() security.Authenticator {
	return d.apiKeyAuthenticator
}
