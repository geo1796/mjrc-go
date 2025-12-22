package admin

import (
	"mjrc/api/routes/admin/middlewares/is_authed"
	"mjrc/api/routes/admin/routes/admin_auth"
	"mjrc/core/chix"
	"mjrc/core/runtime"
)

const Prefix = "/admin"

func Group(deps runtime.Dependencies) *chix.Group {
	g := chix.NewGroup(Prefix)

	g.Add(
		admin_auth.Route(deps.JWT(), deps.APIConfig().AdminPassword),

		is_authed.Middleware(deps.JWT()),
	)

	return g
}
