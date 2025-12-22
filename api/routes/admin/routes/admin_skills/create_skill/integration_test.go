package create_skill

import (
	"bytes"
	"context"
	"encoding/json"
	"mjrc/core/postgres"
	"testing"

	"mjrc/core/models"
	"net/http"
	"net/http/httptest"
	"sort"

	"github.com/go-chi/chi/v5"
)

func TestIntegration_CreateSkill(t *testing.T) {
	ctx := context.Background()
	container, db := postgres.NewTestContainer(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	r := chi.NewRouter()
	Route(db).Register(r)

	// happy path: valid payload -> 201 and row inserted
	in := models.Skill{
		Name:             "double under",
		Level:            5,
		YoutubeVideoID:   "yt_double_under",
		IsVideoLandscape: true,
		Categories:       []models.SkillCategory{models.SkillCategoryMultiples, models.SkillCategoryBasics},
	}

	body, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPost, Path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	// verify it's persisted
	rows, err := db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 skill in db, got %d", len(rows))
	}
	got := rows[0]
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
	// compare categories ignoring order
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

func TestIntegration_CreateSkill_BadCategory(t *testing.T) {
	ctx := context.Background()
	container, db := postgres.NewTestContainer(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	r := chi.NewRouter()
	Route(db).Register(r)

	in := models.Skill{
		Name:             "invalid category skill",
		Level:            3,
		YoutubeVideoID:   "yt_invalid_cat",
		IsVideoLandscape: false,
		Categories:       []models.SkillCategory{"not-a-valid-category"},
	}

	body, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPost, Path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	// ensure nothing was inserted
	rows, err := db.Queries().GetSkills(ctx)
	if err != nil {
		t.Fatalf("GetSkills failed: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("expected 0 skills in db, got %d", len(rows))
	}
}
