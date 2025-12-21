package public

import (
	"mjrc/api/routes/public/get_skills"
	"mjrc/core/chix"
	"mjrc/core/runtime"
)

const Prefix = "/public"

func Group(deps runtime.Dependencies) *chix.Group {
	group := chix.NewGroup(Prefix)

	group.Add(
		get_skills.Route(deps.DB()),
	)

	return group
}
