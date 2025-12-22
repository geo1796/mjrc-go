package create_skill

import (
	"encoding/json"
	"mjrc/core/models"
	"mjrc/core/postgres"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"mjrc/core/postgres/dao"
)

type Handler interface {
	createSkillForAdmin(w http.ResponseWriter, r *http.Request)
}

type handler struct {
    db postgres.DB
}

func (h *handler) createSkillForAdmin(w http.ResponseWriter, r *http.Request) {
	//	junie
	// use models.Skill as input
	// insert using sqlc generated queries in pkg dao

	// Decode JSON body
	var in models.Skill
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Normalize slices (NOT NULL array columns should not receive NULL)
	categories := make([]string, 0, len(in.Categories))
	for _, c := range in.Categories {
		categories = append(categories, string(c))
	}
	prereqs := make([]pgtype.UUID, 0, len(in.Prerequisites))
	for _, p := range in.Prerequisites {
		prereqs = append(prereqs, pgtype.UUID{Bytes: uuid.UUID(p), Valid: true})
	}

	params := dao.CreateSkillParams{
		Name:             in.Name,
		YoutubeVideoID:   in.YoutubeVideoID,
		IsVideoLandscape: in.IsVideoLandscape,
		Level:            in.Level,
		Categories:       categories,
		Prerequisites:    prereqs,
	}

	if err := h.db.Queries().CreateSkill(r.Context(), params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
