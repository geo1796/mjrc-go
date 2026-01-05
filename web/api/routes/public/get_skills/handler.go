package get_skills

import (
	"encoding/json"
	"fmt"
	"mjrc/core/logger"
	"mjrc/core/models"
	"mjrc/core/postgres"
	"net/http"

	"github.com/google/uuid"
)

type Handler interface {
	getSkillsForUser(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	db postgres.DB
}

func (h *handler) getSkillsForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1) Compute an inexpensive fingerprint (count + max(updated_at))
	fp, err := h.db.Queries().SkillsFingerprint(ctx)
	if err != nil {
		logger.Error("failed to query skills fingerprint", logger.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 2) Build an ETag (weak ETag is fine for JSON lists)
	etag := fmt.Sprintf(`W/"skills:%d:%d"`, fp.Cnt, fp.MaxUpdatedAt.Time.Unix())

	// 3) Conditional request
	if inm := r.Header.Get("If-None-Match"); inm != "" && inm == etag {
		w.Header().Set("ETag", etag)
		w.WriteHeader(http.StatusNotModified) // 304
		return
	}

	// 4) Normal response
	w.Header().Set("ETag", etag)

	// Query all skills using the generated sqlc dao via our db wrapper
	rows, err := h.db.Queries().GetSkills(ctx)
	if err != nil {
		logger.Error("failed to query skills", logger.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Map dao rows to public API model
	dto := make([]models.Skill, 0, len(rows))
	for _, s := range rows {
		skill := models.Skill{
			ID:               uuid.UUID(s.ID.Bytes),
			Name:             s.Name,
			Level:            s.Level,
			YoutubeVideoID:   s.YoutubeVideoID,
			IsVideoLandscape: s.IsVideoLandscape,
			Prerequisites:    make([]uuid.UUID, len(s.Prerequisites)),
			Categories:       make([]models.SkillCategory, len(s.Categories)),
			CreatedAt:        s.CreatedAt.Time.UTC(),
			UpdatedAt:        s.UpdatedAt.Time.UTC(),
		}

		for i, c := range s.Categories {
			skill.Categories[i] = models.SkillCategory(c)
		}
		for i, p := range s.Prerequisites {
			skill.Prerequisites[i] = p.Bytes
		}

		dto = append(dto, skill)
	}

	if err = json.NewEncoder(w).Encode(dto); err != nil {
		logger.Error("failed to encode skills", logger.Err(err))
	}
}
