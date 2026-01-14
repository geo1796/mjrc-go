package max_bytes

import (
	"mjrc/web/chix"
	"net/http"
)

const Name = "max_bytes"

// Middleware limits the maximum size of an incoming request body.
// It wraps r.Body with http.MaxBytesReader so handlers will get a 413
// if they try to read more than maxBytes.
//
// Note: this does not *prevent* a handler from being called; it prevents
// reading more than the limit.
func Middleware(maxBytes int64) *chix.Middleware {
	return chix.NewMiddleware(Name, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only wrap when there is a body.
			if r.Body != nil && r.Body != http.NoBody {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}

			next.ServeHTTP(w, r)
		})
	})
}
