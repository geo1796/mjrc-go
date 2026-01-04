package security

type Authenticator interface {
	Authenticate(string) bool
}

type authenticator struct {
	secret string
}

func (p *authenticator) Authenticate(secret string) bool {
	return p.secret == secret
}

func NewAuthenticator(secret string) Authenticator {
	return &authenticator{secret}
}
