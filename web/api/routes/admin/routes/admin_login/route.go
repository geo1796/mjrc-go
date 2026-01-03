package admin_login

import (
	"mjrc/core/runtime"
	"mjrc/web/chix"
	"net/http"
)

const (
	Path   = "/login"
	Method = http.MethodPost
)

func Route(deps runtime.Dependencies) *chix.Route {
	hdlr := &handler{
		jwt:                deps.JWT(),
		adminAuthenticator: deps.AdminAuthenticator(),
	}
	return chix.NewRoute(Path, Method, hdlr.login)
}
