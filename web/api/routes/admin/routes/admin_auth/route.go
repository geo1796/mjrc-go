package admin_auth

import (
	"mjrc/core/security"
	"mjrc/web/chix"
	"net/http"
)

const (
	Path   = "/auth"
	Method = http.MethodPost
)

func Route(jwt security.JWT, adminPassword security.Authenticator) *chix.Route {
	hdlr := &handler{
		jwt:                jwt,
		adminAuthenticator: adminPassword,
	}
	return chix.NewRoute(Path, Method, hdlr.authenticateUser)
}
