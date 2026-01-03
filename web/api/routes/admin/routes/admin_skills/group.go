package admin_skills

import (
	"mjrc/core/runtime"
	"mjrc/web/api/routes/admin/routes/admin_skills/create_skill"
	"mjrc/web/api/routes/admin/routes/admin_skills/delete_skill"
	"mjrc/web/api/routes/admin/routes/admin_skills/update_skill"
	"mjrc/web/chix"
)

const Prefix = "/skills"

func Group(deps runtime.Dependencies) *chix.Group {
	group := chix.NewGroup(Prefix)

	group.Add(
		create_skill.Route(deps),
		update_skill.Route(deps),
		delete_skill.Route(deps),
	)

	return group
}
