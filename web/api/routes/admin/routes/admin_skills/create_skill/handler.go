package create_skill

import (
	"encoding/json"
	"mjrc/core/logger"
	"mjrc/core/models"
	"mjrc/core/postgres"
	"net/http"

	"mjrc/core/postgres/dao"

	"github.com/jackc/pgx/v5/pgtype"
)

type Handler interface {
	createSkillForAdmin(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	db postgres.DB
}

func (h *handler) createSkillForAdmin(w http.ResponseWriter, r *http.Request) {
	// Decode JSON body
	var in models.Skill
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		logger.Error("failed to decode JSON body", logger.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Normalize slices (NOT NULL array columns should not receive NULL)
	categories := make([]string, 0, len(in.Categories))
	for _, c := range in.Categories {
		if !c.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		categories = append(categories, string(c))
	}
	prereqs := make([]pgtype.UUID, 0, len(in.Prerequisites))
	for _, p := range in.Prerequisites {
		prereqs = append(prereqs, pgtype.UUID{Bytes: p, Valid: true})
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
		logger.Error("failed to create skill", logger.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
