package chix

import "github.com/go-chi/chi/v5"

type Component interface {
	Register(router chi.Router)
}
