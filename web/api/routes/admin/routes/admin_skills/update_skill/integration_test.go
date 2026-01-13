package update_skill

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mjrc/core/models"
	"mjrc/core/postgres"
	"mjrc/core/postgres/dao"
	"mjrc/core/runtime"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestIntegration_UpdateSkill(t *testing.T) {
	ctx := context.Background()
	container, db := postgres.NewTestContainer(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	r := chi.NewRouter()
	Route(runtime.NewBuilder().WithDB(db).Build()).Register(r)

	// Seed initial skill
	seed := dao.CreateSkillParams{
		Name:             "double under",
		YoutubeVideoID:   "yt_seed",
		IsVideoLandscape: true,
		Level:            5,
		Categories:       []string{string(models.SkillCategoryMultiples), string(models.SkillCategoryBasics)},
		Prerequisites:    []pgtype.UUID{},
	}
	if err := db.Queries().CreateSkill(ctx, seed); err != nil {
		t.Fatalf("failed to seed skill: %v", err)
	}

	rows, err := db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 skill in db, got %d", len(rows))
	}
	id := rows[0].ID.Bytes

	// Prepare update input (intentionally set a different body ID; handler must use path ID)
	bodyID := uuid.New()
	in := models.Skill{
		ID:               bodyID,
		Name:             "double under (updated)",
		Level:            6,
		YoutubeVideoID:   "yt_updated",
		IsVideoLandscape: false,
		Categories:       []models.SkillCategory{models.SkillCategoryBasics, models.SkillCategoryFootwork},
		Prerequisites:    []uuid.UUID{},
	}

	b, _ := json.Marshal(in)
	url := fmt.Sprintf("/%s", uuid.UUID(id).String())
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	// Verify the row was updated by path ID
	rows, err = db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 skill in db, got %d", len(rows))
	}
	got := rows[0]

	if uuid.UUID(got.ID.Bytes) != uuid.UUID(id) {
		t.Fatalf("expected id to remain the same: %s != %s", uuid.UUID(got.ID.Bytes), uuid.UUID(id))
	}
	if got.Name != in.Name {
		t.Fatalf("expected name %s, got %s", in.Name, got.Name)
	}
	if got.YoutubeVideoID != in.YoutubeVideoID {
		t.Fatalf("expected youtubeVideoId %s, got %s", in.YoutubeVideoID, got.YoutubeVideoID)
	}
	if got.IsVideoLandscape != in.IsVideoLandscape {
		t.Fatalf("expected isVideoLandscape %v, got %v", in.IsVideoLandscape, got.IsVideoLandscape)
	}
	if got.Level != in.Level {
		t.Fatalf("expected level %d, got %d", in.Level, got.Level)
	}

	expCats := []string{string(in.Categories[0]), string(in.Categories[1])}
	sort.Strings(expCats)
	cats := append([]string(nil), got.Categories...)
	sort.Strings(cats)
	if len(cats) != len(expCats) {
		t.Fatalf("expected %d categories, got %d", len(expCats), len(cats))
	}
	for i := range cats {
		if cats[i] != expCats[i] {
			t.Fatalf("expected categories %v, got %v", expCats, got.Categories)
		}
	}
}

func TestIntegration_UpdateSkill_BadCategory(t *testing.T) {
	ctx := context.Background()
	container, db := postgres.NewTestContainer(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	r := chi.NewRouter()
	Route(runtime.NewBuilder().WithDB(db).Build()).Register(r)

	// Seed initial skill
	seed := dao.CreateSkillParams{
		Name:             "cross",
		YoutubeVideoID:   "yt_cross",
		IsVideoLandscape: true,
		Level:            2,
		Categories:       []string{string(models.SkillCategoryBasics)},
		Prerequisites:    []pgtype.UUID{},
	}
	if err := db.Queries().CreateSkill(ctx, seed); err != nil {
		t.Fatalf("failed to seed skill: %v", err)
	}
	rows, err := db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 skill in db, got %d", len(rows))
	}
	id := uuid.UUID(rows[0].ID.Bytes)

	// Prepare invalid category update
	in := models.Skill{
		Name:             "cross updated",
		Level:            3,
		YoutubeVideoID:   "yt_cross_upd",
		IsVideoLandscape: false,
		Categories:       []models.SkillCategory{"not-a-valid-category"},
	}
	b, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPut, "/"+id.String(), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	// Ensure row not updated
	rows, err = db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if rows[0].Name != seed.Name {
		t.Fatalf("expected name unchanged %s, got %s", seed.Name, rows[0].Name)
	}
}

func TestIntegration_UpdateSkill_DuplicateCategories(t *testing.T) {
	ctx := context.Background()
	container, db := postgres.NewTestContainer(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	r := chi.NewRouter()
	Route(runtime.NewBuilder().WithDB(db).Build()).Register(r)

	// Seed initial skill
	seed := dao.CreateSkillParams{
		Name:             "speed step",
		YoutubeVideoID:   "yt_speed_step",
		IsVideoLandscape: true,
		Level:            3,
		Categories:       []string{string(models.SkillCategoryBasics)},
		Prerequisites:    []pgtype.UUID{},
	}
	if err := db.Queries().CreateSkill(ctx, seed); err != nil {
		t.Fatalf("failed to seed skill: %v", err)
	}
	rows, err := db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 skill in db, got %d", len(rows))
	}
	id := uuid.UUID(rows[0].ID.Bytes)

	// Attempt update with duplicate categories
	in := models.Skill{
		Name:             "speed step updated",
		Level:            4,
		YoutubeVideoID:   "yt_speed_step_upd",
		IsVideoLandscape: false,
		Categories:       []models.SkillCategory{models.SkillCategoryBasics, models.SkillCategoryBasics},
	}
	b, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPut, "/"+id.String(), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	// Ensure row not updated
	rows, err = db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if rows[0].Name != seed.Name {
		t.Fatalf("expected name unchanged %s, got %s", seed.Name, rows[0].Name)
	}
}

func TestIntegration_UpdateSkill_DuplicatePrerequisites(t *testing.T) {
	ctx := context.Background()
	container, db := postgres.NewTestContainer(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	r := chi.NewRouter()
	Route(runtime.NewBuilder().WithDB(db).Build()).Register(r)

	// Seed initial skill
	seed := dao.CreateSkillParams{
		Name:             "side swing",
		YoutubeVideoID:   "yt_side_swing",
		IsVideoLandscape: false,
		Level:            2,
		Categories:       []string{string(models.SkillCategoryBasics)},
		Prerequisites:    []pgtype.UUID{},
	}
	if err := db.Queries().CreateSkill(ctx, seed); err != nil {
		t.Fatalf("failed to seed skill: %v", err)
	}
	rows, err := db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 skill in db, got %d", len(rows))
	}
	id := uuid.UUID(rows[0].ID.Bytes)

	dup := uuid.New()
	in := models.Skill{
		Name:             "side swing updated",
		Level:            3,
		YoutubeVideoID:   "yt_side_swing_upd",
		IsVideoLandscape: true,
		Categories:       []models.SkillCategory{models.SkillCategoryBasics},
		Prerequisites:    []uuid.UUID{dup, dup},
	}
	b, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPut, "/"+id.String(), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	rows, err = db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if rows[0].Name != seed.Name {
		t.Fatalf("expected name unchanged %s, got %s", seed.Name, rows[0].Name)
	}
}

func TestIntegration_UpdateSkill_NotFound(t *testing.T) {
	ctx := context.Background()
	container, db := postgres.NewTestContainer(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	r := chi.NewRouter()
	Route(runtime.NewBuilder().WithDB(db).Build()).Register(r)

	in := models.Skill{
		Name:             "nonexistent",
		Level:            1,
		YoutubeVideoID:   "yt_nonexistent",
		IsVideoLandscape: true,
		Categories:       []models.SkillCategory{models.SkillCategoryBasics},
		Prerequisites:    []uuid.UUID{},
	}
	b, _ := json.Marshal(in)

	id := uuid.New()
	req := httptest.NewRequest(http.MethodPut, "/"+id.String(), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestIntegration_UpdateSkill_EmptyCategories(t *testing.T) {
	ctx := context.Background()
	container, db := postgres.NewTestContainer(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	r := chi.NewRouter()
	Route(runtime.NewBuilder().WithDB(db).Build()).Register(r)

	// Seed initial skill
	seed := dao.CreateSkillParams{
		Name:             "basic jump",
		YoutubeVideoID:   "yt_basic_jump",
		IsVideoLandscape: true,
		Level:            1,
		Categories:       []string{string(models.SkillCategoryBasics)},
		Prerequisites:    []pgtype.UUID{},
	}
	if err := db.Queries().CreateSkill(ctx, seed); err != nil {
		t.Fatalf("failed to seed skill: %v", err)
	}
	rows, err := db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 skill in db, got %d", len(rows))
	}
	id := uuid.UUID(rows[0].ID.Bytes)

	// Attempt update with empty categories
	in := models.Skill{
		Name:             "basic jump updated",
		Level:            2,
		YoutubeVideoID:   "yt_basic_jump_upd",
		IsVideoLandscape: false,
		Categories:       []models.SkillCategory{},
	}
	b, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPut, "/"+id.String(), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	// Ensure row not updated (name unchanged and categories still present)
	rows, err = db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if rows[0].Name != seed.Name {
		t.Fatalf("expected name unchanged %s, got %s", seed.Name, rows[0].Name)
	}
	if len(rows[0].Categories) == 0 {
		t.Fatalf("expected categories to remain non-empty")
	}
}
