package chix

import (
	"mjrc/core/logger"

	"github.com/go-chi/chi/v5"
)

type Group struct {
	prefix     string
	components []Component
}

func (g *Group) Register(router chi.Router) {
	subRouter := chi.NewRouter()

	logger.Info("Registering group...",
		logger.Any("prefix", g.prefix),
		logger.Any("components", len(g.components)))

	for _, c := range g.components {
		c.Register(subRouter)
	}

	router.Mount(g.prefix, subRouter)

	logger.Info("Group registered",
		logger.Any("prefix", g.prefix),
		logger.Any("components", len(g.components)))
}

func NewGroup(prefix string) *Group {
	return &Group{prefix: prefix}
}

func (g *Group) Add(components ...Component) {
	g.components = append(g.components, components...)
}
