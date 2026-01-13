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

	if len(in.Categories) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate and normalize categories (reject duplicates)
	seenCats := make(map[string]struct{}, len(in.Categories))
	categories := make([]string, len(in.Categories))
	for i, c := range in.Categories {
		if !c.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cat := string(c)
		if _, seen := seenCats[cat]; seen {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		seenCats[cat] = struct{}{}
		categories[i] = cat
	}

	// Normalize prerequisites (reject duplicates)
	seenPrereqs := make(map[[16]byte]struct{}, len(in.Prerequisites))
	prereqs := make([]pgtype.UUID, len(in.Prerequisites))
	for i, p := range in.Prerequisites {
		if _, seen := seenPrereqs[p]; seen {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		seenPrereqs[p] = struct{}{}
		prereqs[i] = pgtype.UUID{Bytes: p, Valid: true}
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
