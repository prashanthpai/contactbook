package auth

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicAuth(t *testing.T) {

	assert := require.New(t)

	// auth pass case
	passHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	validCfg := &Config{
		User:         "rightuser",
		PasswordHash: sha256.Sum256([]byte("rightpassword")),
	}

	basicAuthPass := BasicAuth(validCfg, passHandler)

	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.SetBasicAuth(validCfg.User, "rightpassword")
	w := httptest.NewRecorder()
	basicAuthPass.ServeHTTP(w, req)
	assert.Equal(w.Code, http.StatusOK)

	// auth fail cases

	failHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Fatal("Auth middleware failed to block unauthorized request")
	})

	inputs := []struct {
		User     string
		Password string
	}{
		{"", ""},
		{"", "rightpassword"},
		{"rightuser", ""},
		{"rightuser", "wrongpassword"},
		{"wronguser", "rightpassword"},
		{"wronguser", "wrongpassword"},
	}

	basicAuthFail := BasicAuth(validCfg, failHandler)

	for _, cfg := range inputs {
		req := httptest.NewRequest("GET", "http://localhost", nil)
		if cfg.User != "" && cfg.Password != "" {
			req.SetBasicAuth(cfg.User, cfg.Password)
		}
		w := httptest.NewRecorder()
		basicAuthFail.ServeHTTP(w, req)
		assert.Equal(w.Body.String(), fmt.Sprintf("%s\n", bodyStr))
		assert.Equal(w.Code, http.StatusUnauthorized)
		assert.Equal(w.Result().Header.Get(authHeaderStr), realmStr)
	}
}
