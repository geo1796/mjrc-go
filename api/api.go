package api

import (
	"mjrc/core/chix"
	"mjrc/core/runtime"

	"github.com/go-chi/chi/v5"
)

const Prefix = "/api"

func group(router chi.Router, deps runtime.Dependencies) *chix.Group {
	group := chix.NewGroup(Prefix)

	group.Add(
	//TODO: add routes
	)

	return group
}

func Register(router chi.Router, deps runtime.Dependencies) {
	group(router, deps).Register(router)
}
