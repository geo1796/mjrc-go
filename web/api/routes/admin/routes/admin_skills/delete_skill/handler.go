package delete_skill

import (
	"mjrc/core/logger"
	"mjrc/core/postgres"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler interface {
	deleteSkillForAdmin(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	db postgres.DB
}

func (h *handler) deleteSkillForAdmin(w http.ResponseWriter, r *http.Request) {
	// Extract and validate UUID from path
	uid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		logger.Error("failed to parse UUID from path", logger.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Perform deletion
	if err = h.db.Queries().DeleteSkill(r.Context(), pgtype.UUID{Bytes: uid, Valid: true}); err != nil {
		logger.Error("failed to delete skill", logger.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// No content on success
	w.WriteHeader(http.StatusNoContent)
}
