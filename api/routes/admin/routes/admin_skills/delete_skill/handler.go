package delete_skill

import (
	"mjrc/core/postgres"
	"net/http"
)

type Handler interface {
	deleteSkillForAdmin(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	db postgres.DB
}

func (h *handler) deleteSkillForAdmin(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
