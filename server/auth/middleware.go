package auth

import (
	"net/http"
)

const (
	realmStr      = `Basic realm="Restricted Contact Book", charset="UTF-8"`
	bodyStr       = "Unauthorized request."
	authHeaderStr = "WWW-Authenticate"
)

type Config struct {
	User     string
	Password string
}

func BasicAuth(config *Config, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set(authHeaderStr, realmStr)
			http.Error(w, bodyStr, http.StatusUnauthorized)
			return
		}

		if user != config.User || password != config.Password {
			w.Header().Set(authHeaderStr, realmStr)
			http.Error(w, bodyStr, http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
