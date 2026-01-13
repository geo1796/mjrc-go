package update_skill

import (
	"encoding/json"
	"errors"
	"mjrc/core/logger"
	"mjrc/core/models"
	"mjrc/core/postgres"
	"mjrc/core/postgres/dao"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	params := dao.UpdateSkillParams{
		ID:               pgtype.UUID{Bytes: parsedID, Valid: true}, // use path id, ignore body id
		Name:             in.Name,
		YoutubeVideoID:   in.YoutubeVideoID,
		IsVideoLandscape: in.IsVideoLandscape,
		Level:            in.Level,
		Categories:       categories,
		Prerequisites:    prereqs,
	}

	if _, err = h.db.Queries().UpdateSkill(r.Context(), params); err != nil {
		logger.Error("failed to update skill", logger.Err(err))

		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
