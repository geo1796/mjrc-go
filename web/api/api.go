package api

import (
	"mjrc/core/runtime"
	"mjrc/web/api/middlewares/api_key"
	"mjrc/web/api/routes/admin"
	"mjrc/web/api/routes/public"
	"mjrc/web/chix"

	"github.com/go-chi/chi/v5"
)

const Prefix = "/api"

func group(deps runtime.Dependencies) *chix.Group {
	group := chix.NewGroup(Prefix, api_key.Middleware(deps))

	group.Add(
		public.Group(deps),
		admin.Group(deps),
	)

	return group
}

func Register(router chi.Router, deps runtime.Dependencies) {
	group(deps).Register(router)
}
