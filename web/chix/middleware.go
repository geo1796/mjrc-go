package chix

import (
	"mjrc/core/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Middleware struct {
	name string
	fn   func(next http.Handler) http.Handler
}

func (m *Middleware) Register(router chi.Router) {
	router.Use(m.fn)
	logger.Info("Middleware registered", logger.Any("name", m.name))
}

func NewMiddleware(name string, fn func(next http.Handler) http.Handler) *Middleware {
	return &Middleware{name, fn}
}
