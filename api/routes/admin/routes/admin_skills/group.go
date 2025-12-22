package admin_skills

import (
	"mjrc/api/routes/admin/routes/admin_skills/create_skill"
	"mjrc/api/routes/admin/routes/admin_skills/delete_skill"
	"mjrc/api/routes/admin/routes/admin_skills/update_skill"
	"mjrc/core/chix"
	"mjrc/core/runtime"
)

const Prefix = "/skills"

func Group(deps runtime.Dependencies) *chix.Group {
	group := chix.NewGroup(Prefix)

	group.Add(
		create_skill.Route(deps.DB()),
		update_skill.Route(deps.DB()),
		delete_skill.Route(deps.DB()),
	)

	return group
}
