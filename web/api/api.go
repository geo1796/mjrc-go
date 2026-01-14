package api

import (
	"mjrc/core/runtime"
	"mjrc/web/api/middlewares/api_key"
	"mjrc/web/api/middlewares/max_bytes"
	"mjrc/web/api/routes/admin"
	"mjrc/web/api/routes/public"
	"mjrc/web/chix"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

const Prefix = "/api"

func group(deps runtime.Dependencies) *chix.Group {
	api := chix.NewGroup(Prefix,
		chix.NewMiddleware("content_type", middleware.SetHeader("Content-Type", "application/json; charset=utf-8")),
		chix.NewMiddleware("rate_limit", httprate.LimitByIP(100, time.Minute)))

	api.Add(
		max_bytes.Middleware(1<<20), // 1 MiB
		api_key.Middleware(deps),
		public.Group(deps),
		admin.Group(deps),
	)

	return api
}

func Register(router chi.Router, deps runtime.Dependencies) {
	group(deps).Register(router)
}
