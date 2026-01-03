package runtime

import (
	"mjrc/core/postgres"
	"mjrc/core/security"
)

type Builder struct {
	deps *dependencies
}

func NewBuilder() *Builder {
	return &Builder{&dependencies{}}
}

func (b *Builder) WithDB(db postgres.DB) *Builder {
	b.deps.db = db
	return b
}

func (b *Builder) WithJWT(jwt security.JWT) *Builder {
	b.deps.jwt = jwt
	return b
}

func (b *Builder) WithAdminAuthenticator(adminAuthenticator security.Authenticator) *Builder {
	b.deps.adminAuthenticator = adminAuthenticator
	return b
}

func (b *Builder) WithAPIKeyAuthenticator(apiKeyAuthenticator security.Authenticator) *Builder {
	b.deps.apiKeyAuthenticator = apiKeyAuthenticator
	return b
}

func (b *Builder) Build() Dependencies {
	deps := b.deps
	b.deps = nil
	return deps
}
