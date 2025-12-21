package admin_skills

import (
	"mjrc/core/chix"

	"github.com/go-chi/chi/v5"
)

const Prefix = "/skills"

func group(router chi.Router) *chix.Group {
	group := chix.NewGroup(Prefix)

	group.Add(
	//TODO: add routes
	)

	return group
}
