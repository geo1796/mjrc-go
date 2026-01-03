package security

type Authenticator interface {
	Authenticate(string) bool
}

type authenticator struct {
	secret string
}

func (p *authenticator) Authenticate(other string) bool {
	return p.secret == other
}

func NewAuthenticator(secret string) Authenticator {
	return &authenticator{secret}
}
