package admin

import (
	"mjrc/core/runtime"
	"mjrc/web/api/routes/admin/middlewares/authenticate_admin"
	"mjrc/web/api/routes/admin/routes/admin_auth"
	"mjrc/web/api/routes/admin/routes/admin_skills"
	"mjrc/web/chix"
)

const Prefix = "/admin"

func Group(deps runtime.Dependencies) *chix.Group {
	g := chix.NewGroup(Prefix)

	g.Add(
		admin_auth.Route(deps.JWT(), deps.AdminAuthenticator()),

		authenticate_admin.Middleware(deps.JWT()),

		admin_skills.Group(deps),
	)

	return g
}
