package public

import (
	"mjrc/core/runtime"
	"mjrc/web/api/routes/public/get_skills"
	"mjrc/web/chix"
)

const Prefix = "/"

func Group(deps runtime.Dependencies) *chix.Group {
	return chix.NewGroup(Prefix,
		get_skills.Route(deps),
	)
}
