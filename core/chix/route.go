package chix

import (
	"mjrc/core/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	path, method string
	fn           http.HandlerFunc
}

func (r *Route) Register(router chi.Router) {
	router.Method(r.method, r.path, r.fn)
	logger.Info("Route registered",
		logger.Any("path", r.path),
		logger.Any("method", r.method))
}

func NewRoute(path, method string, fn http.HandlerFunc) *Route {
	return &Route{path, method, fn}
}
