package admin_login

import (
	"encoding/json"
	"mjrc/core/logger"
	"mjrc/core/security"
	"net/http"
)

type Handler interface {
	login(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	jwt                security.JWT
	adminAuthenticator security.Authenticator
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var i input
	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		logger.Error("failed to decode JSON body", logger.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !h.adminAuthenticator.Authenticate(i.Password) {
		logger.Warn("invalid password", logger.Any("password", i.Password))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, _, err := h.jwt.Generate()
	if err != nil {
		logger.Error("failed to generate JWT", logger.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(output{token}); err != nil {
		logger.Error("failed to encode JWT", logger.Err(err))
	}
}
