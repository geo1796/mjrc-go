package admin

import (
	"mjrc/core/runtime"
	"mjrc/web/api/routes/admin/middlewares/authenticate_admin"
	"mjrc/web/api/routes/admin/routes/admin_login"
	"mjrc/web/api/routes/admin/routes/admin_skills"
	"mjrc/web/chix"
)

const Prefix = "/admin"

func Group(deps runtime.Dependencies) *chix.Group {
	g := chix.NewGroup(Prefix, admin_login.Route(deps))

	g.Add(chix.NewGroup("/", authenticate_admin.Middleware(deps),
		admin_skills.Group(deps),
	))

	return g
}
