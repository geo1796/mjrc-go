package public

import (
	"mjrc/core/runtime"
	"mjrc/web/api/routes/public/get_skills"
	"mjrc/web/chix"
)

const Prefix = "/public"

func Group(deps runtime.Dependencies) *chix.Group {
	group := chix.NewGroup(Prefix)

	group.Add(
		get_skills.Route(deps),
	)

	return group
}
