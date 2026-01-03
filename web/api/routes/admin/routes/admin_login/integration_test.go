package admin_login

import (
	"bytes"
	"encoding/json"
	"mjrc/core/runtime"
	"mjrc/core/security"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

type loginReq struct {
	Password string `json:"password"`
}

func TestIntegration_AuthenticateUser(t *testing.T) {
	adminPassword := "Test123!"

	jwt := security.NewJWT(
		"jwt",
		[]byte("secret"),
		1*time.Minute,
	)

	r := chi.NewRouter()
	Route(runtime.NewBuilder().
		WithAdminAuthenticator(security.NewAuthenticator("Test123!")).
		WithJWT(jwt).Build()).
		Register(r)

	t.Run("bad json -> 400 and no cookie", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, Path, bytes.NewBufferString("{bad json}"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		for _, c := range rec.Result().Cookies() {
			if c.Name == jwt.CookieName() {
				t.Fatalf("should not set jwt cookie on bad json")
			}
		}
	})

	t.Run("wrong password -> 401 and no cookie", func(t *testing.T) {
		body, _ := json.Marshal(loginReq{Password: "wrong"})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, Path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
		for _, c := range rec.Result().Cookies() {
			if c.Name == jwt.CookieName() {
				t.Fatalf("should not set jwt cookie on unauthorized")
			}
		}
	})

	t.Run("correct password -> 200, jwt cookie set and valid", func(t *testing.T) {
		body, _ := json.Marshal(loginReq{Password: adminPassword})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, Path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var cookie *http.Cookie
		for _, c := range rec.Result().Cookies() {
			if c.Name == jwt.CookieName() {
				cookie = c
				break
			}
		}
		if cookie == nil {
			t.Fatalf("expected jwt cookie to be set")
		}

		if !cookie.HttpOnly {
			t.Fatalf("expected HttpOnly cookie")
		}
		if !cookie.Secure {
			t.Fatalf("expected Secure cookie")
		}
		if cookie.Path != "/" {
			t.Fatalf("expected cookie Path=/, got %s", cookie.Path)
		}
		if cookie.SameSite != http.SameSiteStrictMode {
			t.Fatalf("expected SameSite Strict mode, got %v", cookie.SameSite)
		}
		if cookie.MaxAge <= 0 {
			t.Fatalf("expected positive MaxAge, got %d", cookie.MaxAge)
		}
		if time.Now().After(cookie.Expires) {
			t.Fatalf("expected cookie to expire in the future")
		}

		// token should be parseable by the same jwt
		if err := jwt.Parse(cookie.Value); err != nil {
			t.Fatalf("expected valid jwt in cookie, got parse error: %v", err)
		}
	})
}
