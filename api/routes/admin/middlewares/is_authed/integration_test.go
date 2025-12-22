package is_authed

import (
	"mjrc/core/security"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestIntegration_IsAuthed(t *testing.T) {
	jwt := security.NewJWT(
		"jwt",
		[]byte("secret"),
		1*time.Minute,
	)

	protectedPath := "/protected"

	// Router with middleware and a protected route
	r := chi.NewRouter()
	Middleware(jwt).Register(r)

	called := false
	r.Get(protectedPath, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	t.Run("no cookie -> 401", func(t *testing.T) {
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

	t.Run("invalid token -> 401", func(t *testing.T) {
		called = false
		req := httptest.NewRequest(http.MethodGet, protectedPath, nil)
		req.AddCookie(&http.Cookie{Name: jwt.CookieName(), Value: "not-a-valid-token"})
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
		if called {
			t.Fatalf("handler should not be called when unauthorized")
		}
	})

	t.Run("valid token -> 200", func(t *testing.T) {
		called = false
		token, _, err := jwt.Generate()
		if err != nil {
			t.Fatalf("failed generating token: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, protectedPath, nil)
		req.AddCookie(&http.Cookie{Name: jwt.CookieName(), Value: token})
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
