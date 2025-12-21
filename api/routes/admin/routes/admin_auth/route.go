package admin_auth

import (
	"mjrc/core/chix"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	Path   = "/auth"
	Method = http.MethodPost
)

func Route(router chi.Router, password string) *chix.Route {
	panic("not implemented")
}
