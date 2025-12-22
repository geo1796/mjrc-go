package api_key

import "net/http"

type Handler interface {
	validateApiKey(next http.Handler) http.Handler
}

type handler struct {
	apiKey string
}

func (h *handler) validateApiKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := r.Header.Get("X-API-KEY")
		if input != h.apiKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// API key is valid; proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
