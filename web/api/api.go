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
	return chix.NewGroup(Prefix,
		api_key.Middleware(deps),
		public.Group(deps),
		admin.Group(deps),
	)
}

func Register(router chi.Router, deps runtime.Dependencies) {
	group(deps).Register(router)
}
