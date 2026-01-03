package authenticate_admin

import (
	"mjrc/core/runtime"
	"mjrc/core/security"
	"mjrc/web/chix"
	"net/http"
)

const Name = "authenticate_admin"

func Middleware(deps runtime.Dependencies) *chix.Middleware {
	hdlr := &handler{deps.JWT()}
	return chix.NewMiddleware(
		Name,
		hdlr.authenticateAdmin,
	)
}

type Handler interface {
	authenticateAdmin(next http.Handler) http.Handler
}

type handler struct {
	jwt security.JWT
}

func (h *handler) authenticateAdmin(next http.Handler) http.Handler {
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
