package admin_auth

import (
	"encoding/json"
	"mjrc/core/security"
	"net/http"
	"time"
)

type Handler interface {
	authenticateUser(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	adminPassword string
	jwt           security.JWT
}

func (h *handler) authenticateUser(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Password string `json:"password"`
	}

	var i input
	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if i.Password != h.adminPassword {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwt, expiry, err := h.jwt.Generate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.jwt.CookieName(),
		Value:    jwt,
		Expires:  expiry,
		MaxAge:   int(time.Until(expiry).Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
