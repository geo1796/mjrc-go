package get_skills

import (
	"encoding/json"
	"fmt"
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
	skillsRows, err := h.db.Queries().GetSkills(ctx)
	if err != nil {
		// As per project error handling, return only status code on failure
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Map dao rows to public API model
	resp := make([]models.Skill, 0, len(skillsRows))
	for _, s := range skillsRows {
		skill := models.Skill{
			ID:               uuid.UUID(s.ID.Bytes),
			Name:             s.Name,
			Level:            s.Level,
			YoutubeVideoID:   s.YoutubeVideoID,
			IsVideoLandscape: s.IsVideoLandscape,
			Prerequisites:    make([]uuid.UUID, 0, len(s.Prerequisites)),
			Categories:       make([]models.SkillCategory, 0, len(s.Categories)),
			CreatedAt:        s.CreatedAt.Time,
			UpdatedAt:        s.UpdatedAt.Time,
		}

		for _, c := range s.Categories {
			skill.Categories = append(skill.Categories, models.SkillCategory(c))
		}
		for _, p := range s.Prerequisites {
			skill.Prerequisites = append(skill.Prerequisites, uuid.UUID(p.Bytes))
		}

		resp = append(resp, skill)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
