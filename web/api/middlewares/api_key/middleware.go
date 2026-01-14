package api_key

import (
	"mjrc/core/logger"
	"mjrc/core/runtime"
	"mjrc/core/security"
	"mjrc/web/chix"
	"net/http"
)

const Name = "api_key"

func Middleware(deps runtime.Dependencies) *chix.Middleware {
	hdlr := &handler{deps.APIKeyAuthenticator()}
	return chix.NewMiddleware(Name, hdlr.validateApiKey)
}

type Handler interface {
	validateApiKey(next http.Handler) http.Handler
}

type handler struct {
	authenticator security.Authenticator
}

func (h *handler) validateApiKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := r.Header.Get("X-API-KEY")
		if err := h.authenticator.Authenticate(input); err != nil {
			logger.Warn("invalid API key", logger.Err(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// API key is valid; proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
