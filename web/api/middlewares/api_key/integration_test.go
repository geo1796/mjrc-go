package api_key

import (
	"mjrc/core/runtime"
	"mjrc/core/security"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestIntegration_ApiKey(t *testing.T) {
	const protectedPath = "/protected"
	const key = "super-secret-key"

	r := chi.NewRouter()
	deps := runtime.NewBuilder().
		WithAPIKeyAuthenticator(security.NewAPIKeyAuthenticator(key)).
		Build()
	Middleware(deps).Register(r)

	called := false
	r.Get(protectedPath, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	t.Run("no header -> 401", func(t *testing.T) {
		called = false
		req := httptest.NewRequest(http.MethodGet, protectedPath, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
		if called {
			t.Fatalf("handler should not be called when unauthorized")
		}
	})

	t.Run("wrong header -> 401", func(t *testing.T) {
		called = false
		req := httptest.NewRequest(http.MethodGet, protectedPath, nil)
		req.Header.Set("X-API-KEY", "wrong")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
		if called {
			t.Fatalf("handler should not be called when unauthorized")
		}
	})

	t.Run("correct header -> 200", func(t *testing.T) {
		called = false
		req := httptest.NewRequest(http.MethodGet, protectedPath, nil)
		req.Header.Set("X-API-KEY", key)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		if !called {
			t.Fatalf("handler should be called when authorized")
		}
	})
}
