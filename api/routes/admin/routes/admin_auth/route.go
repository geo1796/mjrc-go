package admin_auth

import (
	"mjrc/core/chix"
	"mjrc/core/security"
	"net/http"
)

const (
	Path   = "/auth"
	Method = http.MethodPost
)

func Route(jwt security.JWT, password string) *chix.Route {
	hdlr := &handler{
		jwt:           jwt,
		adminPassword: password,
	}
	return chix.NewRoute(Path, Method, hdlr.authenticateUser)
}
