package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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
		User:     "rightuser",
		Password: "rightpassword",
	}

	basicAuthPass := BasicAuth(validCfg, passHandler)

	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.SetBasicAuth(validCfg.User, validCfg.Password)
	w := httptest.NewRecorder()
	basicAuthPass.ServeHTTP(w, req)
	assert.Equal(w.Code, http.StatusOK)

	// auth fail cases

	failHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Fatal("Auth middleware failed to block unauthorized request")
	})

	inputs := []*Config{
		&Config{"", ""},
		&Config{"", "rightpassword"},
		&Config{"rightuser", ""},
		&Config{"rightuser", "wrongpassword"},
		&Config{"wronguser", "rightpassword"},
		&Config{"wronguser", "wrongpassword"},
	}

	basicAuthFail := BasicAuth(validCfg, failHandler)

	for i, cfg := range inputs {
		req := httptest.NewRequest("GET", "http://localhost", nil)
		req.RemoteAddr = "203.0.113." + strconv.Itoa(i+1) + ":1234"
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

func TestBasicAuthLockoutAfterFailures(t *testing.T) {
	assert := require.New(t)

	validCfg := &Config{
		User:     "rightuser",
		Password: "rightpassword",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	basicAuth := BasicAuth(validCfg, handler)
	clientAddr := "198.51.100.10:9000"

	for i := 0; i < lockoutThreshold; i++ {
		req := httptest.NewRequest("GET", "http://localhost", nil)
		req.RemoteAddr = clientAddr
		req.SetBasicAuth(validCfg.User, "wrongpassword")
		w := httptest.NewRecorder()
		basicAuth.ServeHTTP(w, req)
		assert.Equal(http.StatusUnauthorized, w.Code)
	}

	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.RemoteAddr = clientAddr
	req.SetBasicAuth(validCfg.User, validCfg.Password)
	w := httptest.NewRecorder()
	basicAuth.ServeHTTP(w, req)
	assert.Equal(http.StatusTooManyRequests, w.Code)
	assert.Equal(fmt.Sprintf("%s\n", rateBodyStr), w.Body.String())
}

func TestBasicAuthRateLimitedPerIP(t *testing.T) {
	assert := require.New(t)

	validCfg := &Config{
		User:     "rightuser",
		Password: "rightpassword",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	basicAuth := BasicAuth(validCfg, handler)
	clientAddr := "203.0.113.77:9000"

	for i := 0; i < int(rateLimitBurst); i++ {
		req := httptest.NewRequest("GET", "http://localhost", nil)
		req.RemoteAddr = clientAddr
		req.SetBasicAuth(validCfg.User, validCfg.Password)
		w := httptest.NewRecorder()
		basicAuth.ServeHTTP(w, req)
		assert.Equal(http.StatusOK, w.Code)
	}

	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.RemoteAddr = clientAddr
	req.SetBasicAuth(validCfg.User, validCfg.Password)
	w := httptest.NewRecorder()
	basicAuth.ServeHTTP(w, req)
	assert.Equal(http.StatusTooManyRequests, w.Code)
	assert.Equal(fmt.Sprintf("%s\n", rateBodyStr), w.Body.String())
}
