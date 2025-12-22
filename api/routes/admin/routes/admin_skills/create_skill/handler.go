package create_skill

import (
	"mjrc/core/postgres"
	"net/http"
)

type Handler interface {
	createSkillForAdmin(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	db postgres.DB
}

func (h *handler) createSkillForAdmin(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
	//	junie
	// use models.Skill as input
	// insert using sqlc generated queries in pkg dao
}
