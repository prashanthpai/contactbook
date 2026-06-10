package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

const (
	realmStr      = `Basic realm="Restricted Contact Book", charset="UTF-8"`
	bodyStr       = "Unauthorized request."
	authHeaderStr = "WWW-Authenticate"
)

type Config struct {
	User         string
	PasswordHash [sha256.Size]byte
}

func BasicAuth(config *Config, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set(authHeaderStr, realmStr)
			http.Error(w, bodyStr, http.StatusUnauthorized)
			return
		}

		passwordHash := sha256.Sum256([]byte(password))
		if subtle.ConstantTimeCompare([]byte(user), []byte(config.User)) != 1 ||
			subtle.ConstantTimeCompare(passwordHash[:], config.PasswordHash[:]) != 1 {
			w.Header().Set(authHeaderStr, realmStr)
			http.Error(w, bodyStr, http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
