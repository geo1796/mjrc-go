package admin

import (
	"mjrc/api/routes/admin/middlewares/is_authed"
	"mjrc/api/routes/admin/routes/admin_auth"
	"mjrc/api/routes/admin/routes/admin_skills"
	"mjrc/core/chix"
	"mjrc/core/runtime"
)

const Prefix = "/admin"

func Group(deps runtime.Dependencies) *chix.Group {
	g := chix.NewGroup(Prefix)

	g.Add(
		admin_auth.Route(deps.JWT(), deps.AdminPassword()),

		is_authed.Middleware(deps.JWT()),

		admin_skills.Group(deps),
	)

	return g
}
