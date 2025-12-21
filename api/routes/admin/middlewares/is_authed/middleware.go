package is_authed

import (
	"mjrc/core/chix"

	"github.com/go-chi/chi/v5"
)

const Name = "is_authed"

func Middleware(router chi.Router, jwtSecret string) *chix.Middleware {
	panic("not implemented")
}
