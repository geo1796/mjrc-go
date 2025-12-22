package is_authed

import (
	"mjrc/core/chix"
	"mjrc/core/security"
)

const Name = "is_authed"

func Middleware(jwt security.JWT) *chix.Middleware {
	hdlr := &handler{jwt}
	return chix.NewMiddleware(
		Name,
		hdlr.authenticateRequest,
	)
}
