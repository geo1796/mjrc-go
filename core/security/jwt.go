package security

import (
	"errors"
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type JWT interface {
	Generate() (string, time.Time, error)
	Parse(string) error
}

func NewJWT(secret []byte, ttl time.Duration) JWT {
	return &jwt{
		secret: secret,
		ttl:    ttl,
	}
}

type jwt struct {
	secret []byte
	ttl    time.Duration
}

const issuer = "mjrc-api"
const audience = "mjrc.auth"

func (j *jwt) Generate() (string, time.Time, error) {
	now := time.Now().UTC()
	expiry := now.Add(j.ttl)

	claims := gojwt.RegisteredClaims{
		ExpiresAt: gojwt.NewNumericDate(expiry),
		IssuedAt:  gojwt.NewNumericDate(now),
		Issuer:    issuer,
		Audience:  []string{audience},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)

	if signedToken, err := token.SignedString(j.secret); err != nil {
		return "", time.Time{}, err
	} else {
		return signedToken, expiry, nil
	}
}

func (j *jwt) Parse(raw string) error {
	token, err := gojwt.ParseWithClaims(raw, &gojwt.RegisteredClaims{}, func(token *gojwt.Token) (any, error) {
		if token.Method != gojwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*gojwt.RegisteredClaims)
	if !ok {
		return errors.New("invalid claims type")
	}
	if claims.Issuer != issuer {
		return errors.New("invalid claims issuer")
	}
	if len(claims.Audience) != 1 || claims.Audience[0] != audience {
		return errors.New("invalid claims audience")
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}
