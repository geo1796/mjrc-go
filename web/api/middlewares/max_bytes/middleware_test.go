package max_bytes

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestMiddleware_LimitsRequestBodySize(t *testing.T) {
	const limit int64 = 8

	r := chi.NewRouter()
	// Register the middleware on the router (avoids relying on internal fields of chix.Middleware).
	Middleware(limit).Register(r)

	// Handler reads the whole body. If it exceeds the limit, write 413.
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		var mbErr *http.MaxBytesError
		if errors.As(err, &mbErr) {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	body := bytes.Repeat([]byte("a"), int(limit)+1)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d, got %d", http.StatusRequestEntityTooLarge, rr.Code)
	}
}
