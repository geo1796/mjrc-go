package is_authed

import (
	"mjrc/core/security"
	"net/http"
)

type Handler interface {
	authenticateRequest(next http.Handler) http.Handler
}

type handler struct {
	jwt security.JWT
}

func (h *handler) authenticateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(h.jwt.CookieName())

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err = h.jwt.Parse(cookie.Value); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
