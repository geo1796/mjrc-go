package update_skill

import (
	"mjrc/core/postgres"
	"net/http"
)

type Handler interface {
	updateSkillForAdmin(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	db postgres.DB
}

func (h *handler) updateSkillForAdmin(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
