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

func Route(jwt security.JWT, adminPassword security.AdminPassword) *chix.Route {
	hdlr := &handler{
		jwt:           jwt,
		adminPassword: adminPassword,
	}
	return chix.NewRoute(Path, Method, hdlr.authenticateUser)
}
