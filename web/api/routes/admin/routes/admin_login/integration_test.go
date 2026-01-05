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

func TestIntegration_AuthenticateUser(t *testing.T) {
	adminPassword := "Test123!"

	jwt := security.NewJWT(
		[]byte("secret"),
		1*time.Minute,
	)

	r := chi.NewRouter()
	Route(runtime.NewBuilder().
		WithAdminAuthenticator(security.NewAuthenticator("Test123!")).
		WithJWT(jwt).Build()).
		Register(r)

	t.Run("bad json -> 400 and no Authorization header", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, Path, bytes.NewBufferString("{bad json}"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		if h := rec.Header().Get("Authorization"); h != "" {
			t.Fatalf("should not set Authorization header on bad json")
		}
	})

	t.Run("wrong password -> 401 and no Authorization header", func(t *testing.T) {
		body, _ := json.Marshal(input{Password: "wrong"})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, Path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
		if h := rec.Header().Get("Authorization"); h != "" {
			t.Fatalf("should not set Authorization header on unauthorized")
		}
	})

	t.Run("correct password -> 200, Authorization Bearer set and valid", func(t *testing.T) {
		body, _ := json.Marshal(input{Password: adminPassword})
		before := time.Now().UTC()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, Path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)
		after := time.Now().UTC()

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var got output
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		token := got.Token
		if token == "" {
			t.Fatalf("expected non-empty bearer token")
		}
		// token should be parseable by the same jwt
		if err := jwt.Parse(token); err != nil {
			t.Fatalf("expected valid jwt in Authorization header, got parse error: %v", err)
		}

		// Assert expiry is set correctly based on the TTL (1 minute)
		lower := before.Add(1 * time.Minute)
		upper := after.Add(1 * time.Minute)
		if got.Expiry.Before(lower) || got.Expiry.After(upper) {
			t.Fatalf("expected expiry to be between %v and %v, got %v", lower, upper, got.Expiry)
		}

		// Ensure no cookies set
		if len(rec.Result().Cookies()) > 0 {
			t.Fatalf("expected no cookies to be set")
		}
	})
}
