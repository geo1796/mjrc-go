package get_skills

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/testcontainers/testcontainers-go"

	"mjrc/core/models"
	"mjrc/core/postgres"
	"mjrc/core/postgres/dao"
)

type testData struct {
	skills []models.Skill
}

func buildTestData() (data testData) {
	data.skills = []models.Skill{
		{Name: "skill1", Level: 10, Categories: []models.SkillCategory{"floaters", "multiples"}},
		{Name: "skill2", Level: 5, Categories: []models.SkillCategory{"multiples"}},
		{Name: "skill3", Level: 9, Categories: []models.SkillCategory{"floaters"}},
		{Name: "skill4", Level: 1, Categories: []models.SkillCategory{"floaters"}},
	}
	return
}

func loadTestData(ctx context.Context, t *testing.T) (container testcontainers.Container, pgxProvider postgres.DB, data testData) {
	container, pgxProvider = postgres.NewTestContainer(ctx, t)
	data = buildTestData()

	var err error
	defer func() {
		if err != nil {
			postgres.CleanUpTestContainer(ctx, t, container, pgxProvider)
			t.Fatal(err)
		}
	}()

	q := pgxProvider.Queries()
	for i, s := range data.skills {
		cats := make([]string, 0, len(s.Categories))
		for _, c := range s.Categories {
			cats = append(cats, string(c))
		}
		params := dao.CreateSkillParams{
			Name:             s.Name,
			YoutubeVideoID:   "vid" + s.Name,
			IsVideoLandscape: i%2 == 0,
			Level:            s.Level,
			Categories:       cats,
			Prerequisites:    []pgtype.UUID{},
		}
		err = q.CreateSkill(ctx, params)
		if err != nil {
			return
		}
	}

	return
}

func TestIntegration_ETAG(t *testing.T) {
	ctx := context.Background()
	container, db, _ := loadTestData(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	r := chi.NewRouter()
	Route(db).Register(r)

	// First request to get the ETag
	req1 := httptest.NewRequest(http.MethodGet, Path, nil)
	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("expected status 200 on first request, got %d", rec1.Code)
	}
	etag := rec1.Header().Get("ETag")
	if etag == "" {
		t.Fatalf("expected ETag header on first response")
	}

	// Second request with If-None-Match should yield 304 and no body
	req2 := httptest.NewRequest(http.MethodGet, Path, nil)
	req2.Header.Set("If-None-Match", etag)
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusNotModified {
		t.Fatalf("expected status 304 when ETag matches, got %d", rec2.Code)
	}
	if rec2.Header().Get("ETag") != etag {
		t.Fatalf("expected same ETag to be returned on 304 response")
	}
	if rec2.Body.Len() != 0 {
		t.Fatalf("expected empty body on 304 response, got length %d", rec2.Body.Len())
	}
}

func TestIntegration_GetSkills(t *testing.T) {
	ctx := context.Background()
	container, db, data := loadTestData(ctx, t)
	defer postgres.CleanUpTestContainer(ctx, t, container, db)

	r := chi.NewRouter()
	Route(db).Register(r)

	req := httptest.NewRequest(http.MethodGet, Path, nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var got []models.Skill
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(got) != len(data.skills) {
		t.Fatalf("expected %d skills, got %d", len(data.skills), len(got))
	}

	// Compare ignoring order
	expect := make(map[string]models.Skill)
	for _, s := range data.skills {
		expect[s.Name] = s
	}
	for i := range got {
		gs := got[i]
		es, ok := expect[gs.Name]
		if !ok {
			t.Fatalf("unexpected skill in response: %s", gs.Name)
		}
		if gs.Level != es.Level {
			t.Fatalf("skill %s: expected level %d, got %d", gs.Name, es.Level, gs.Level)
		}
		// compare categories as sets
		if !sameStringSet(catStrings(gs.Categories), catStrings(es.Categories)) {
			t.Fatalf("skill %s: categories mismatch. expected %v, got %v", gs.Name, es.Categories, gs.Categories)
		}
	}
}

func catStrings(cats []models.SkillCategory) []string {
	out := make([]string, len(cats))
	for i, c := range cats {
		out[i] = string(c)
	}
	return out
}

func sameStringSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
