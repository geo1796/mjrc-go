package security

import (
	"crypto/subtle"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Authenticator interface {
	Authenticate(string) error
}

var ErrInvalidSecret = errors.New("invalid secret")

// constantTimeEquals compares two strings in constant time (when lengths match).
func constantTimeEquals(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

type apiKeyAuthenticator struct {
	secret string
}

func (p *apiKeyAuthenticator) Authenticate(secret string) error {
	if !constantTimeEquals(p.secret, secret) {
		return ErrInvalidSecret
	}
	return nil
}

// NewAPIKeyAuthenticator is kept for call-site clarity, but uses the same implementation.
func NewAPIKeyAuthenticator(secret string) Authenticator {
	return &apiKeyAuthenticator{secret}
}

type adminAuthenticator struct {
	hashedPassword string
}

func (p *adminAuthenticator) Authenticate(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashedPassword), []byte(password))
	if err == nil {
		return nil
	}

	// Wrong password → uniform error for callers (avoid leaking details).
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrInvalidSecret
	}

	// Invalid hash / config issue → let caller log/handle as an internal problem.
	return fmt.Errorf("bcrypt compare failed: %w", err)
}

func NewAdminAuthenticator(pwdHash string) Authenticator {
	return &adminAuthenticator{pwdHash}
}
