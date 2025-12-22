package delete_skill

import (
    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgtype"
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
    // Extract and validate UUID from path
    idStr := chi.URLParam(r, "id")
    uid, err := uuid.Parse(idStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // Perform deletion
    if err := h.db.Queries().DeleteSkill(r.Context(), pgtype.UUID{Bytes: uid, Valid: true}); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // No content on success
    w.WriteHeader(http.StatusNoContent)
}
