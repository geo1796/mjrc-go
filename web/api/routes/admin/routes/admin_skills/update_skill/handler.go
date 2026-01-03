package update_skill

import (
	"encoding/json"
	"mjrc/core/logger"
	"mjrc/core/models"
	"mjrc/core/postgres"
	"mjrc/core/postgres/dao"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler interface {
	updateSkillForAdmin(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	db postgres.DB
}

func (h *handler) updateSkillForAdmin(w http.ResponseWriter, r *http.Request) {
	// Extract and validate path id
	parsedID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		logger.Error("failed to parse UUID from path", logger.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Decode JSON body
	var in models.Skill
	if err = json.NewDecoder(r.Body).Decode(&in); err != nil {
		logger.Error("failed to decode JSON body", logger.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate and normalize categories
	categories := make([]string, 0, len(in.Categories))
	for _, c := range in.Categories {
		if !c.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		categories = append(categories, string(c))
	}

	// Normalize prerequisites
	prereqs := make([]pgtype.UUID, 0, len(in.Prerequisites))
	for _, p := range in.Prerequisites {
		prereqs = append(prereqs, pgtype.UUID{Bytes: p, Valid: true})
	}

	params := dao.UpdateSkillParams{
		ID:               pgtype.UUID{Bytes: parsedID, Valid: true}, // use path id, ignore body id
		Name:             in.Name,
		YoutubeVideoID:   in.YoutubeVideoID,
		IsVideoLandscape: in.IsVideoLandscape,
		Level:            in.Level,
		Categories:       categories,
		Prerequisites:    prereqs,
	}

	if err = h.db.Queries().UpdateSkill(r.Context(), params); err != nil {
		logger.Error("failed to update skill", logger.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
