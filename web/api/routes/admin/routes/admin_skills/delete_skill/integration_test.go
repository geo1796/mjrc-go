package delete_skill

import (
	"context"
	"mjrc/core/runtime"
	"net/http"
	"net/http/httptest"
	"testing"

	"mjrc/core/postgres"
	"mjrc/core/postgres/dao"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestIntegration_DeleteSkill(t *testing.T) {
	ctx := context.Background()
	container, db := postgres.NewTestContainer(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	// Seed one skill to delete
	if err := seedOneSkill(ctx, db); err != nil {
		t.Fatalf("failed to seed skill: %v", err)
	}

	rows, err := db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(rows))
	}

	idUUID, _ := uuid.FromBytes(rows[0].ID.Bytes[:])

	r := chi.NewRouter()
	Route(runtime.NewBuilder().WithDB(db).Build()).Register(r)

	req := httptest.NewRequest(http.MethodDelete, "/"+idUUID.String(), nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	// Ensure it's deleted
	rows, err = db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills after delete failed: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("expected 0 skills after delete, got %d", len(rows))
	}
}

func TestIntegration_DeleteSkill_InvalidID(t *testing.T) {
	ctx := context.Background()
	container, db := postgres.NewTestContainer(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	r := chi.NewRouter()
	Route(runtime.NewBuilder().WithDB(db).Build()).Register(r)

	req := httptest.NewRequest(http.MethodDelete, "/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

// seedOneSkill inserts a minimal valid skill row
func seedOneSkill(ctx context.Context, db postgres.DB) error {
	params := dao.CreateSkillParams{
		Name:             "to-delete",
		YoutubeVideoID:   "yt_to_delete",
		IsVideoLandscape: false,
		Level:            3,
		Categories:       []string{"basics"},
		Prerequisites:    []pgtype.UUID{},
	}
	return db.Queries().CreateSkill(ctx, params)
}
