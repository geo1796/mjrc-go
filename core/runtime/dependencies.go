package runtime

import (
	"mjrc/core/postgres"
	"mjrc/core/security"
)

type Dependencies interface {
	DB() postgres.DB
	JWT() security.JWT
	AdminPassword() security.AdminPassword
}

func New(db postgres.DB, jwt security.JWT, adminPassword security.AdminPassword) Dependencies {
	return &dependencies{
		db:            db,
		jwt:           jwt,
		adminPassword: adminPassword,
	}
}

type dependencies struct {
	db            postgres.DB
	jwt           security.JWT
	adminPassword security.AdminPassword
}

func (d *dependencies) DB() postgres.DB {
	return d.db
}

func (d *dependencies) JWT() security.JWT {
	return d.jwt
}

func (d *dependencies) AdminPassword() security.AdminPassword {
	return d.adminPassword
}
