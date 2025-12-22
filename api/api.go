package api

import (
	"mjrc/api/middlewares/api_key"
	"mjrc/api/routes/admin"
	"mjrc/api/routes/public"
	"mjrc/core/chix"
	"mjrc/core/runtime"

	"github.com/go-chi/chi/v5"
)

const Prefix = "/api"

func group(deps runtime.Dependencies, apiKey string) *chix.Group {
	group := chix.NewGroup(Prefix)

	group.Add(
		api_key.Middleware(apiKey),

		public.Group(deps),
		admin.Group(deps),
	)

	return group
}

func Register(router chi.Router, deps runtime.Dependencies, apiKey string) {
	group(deps, apiKey).Register(router)
}
