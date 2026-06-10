package auth

import (
	"log"
	"math"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	realmStr      = `Basic realm="Restricted Contact Book", charset="UTF-8"`
	bodyStr       = "Unauthorized request."
	authHeaderStr = "WWW-Authenticate"
	rateBodyStr   = "Too many authentication attempts."

	rateLimitPerSecond = 5.0
	rateLimitBurst     = 10.0
	lockoutThreshold   = 5
	lockoutDuration    = 30 * time.Second
)

type Config struct {
	User     string
	Password string
}

func BasicAuth(config *Config, h http.Handler) http.Handler {
	guard := newIPGuard()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := sourceIP(r)

		if !guard.allow(ip) {
			http.Error(w, rateBodyStr, http.StatusTooManyRequests)
			return
		}

		if guard.locked(ip) {
			log.Printf("basic auth locked out for ip=%s", ip)
			http.Error(w, rateBodyStr, http.StatusTooManyRequests)
			return
		}

		user, password, ok := r.BasicAuth()
		if !ok {
			attempts := guard.fail(ip)
			log.Printf("basic auth failed from ip=%s user=%q attempts=%d", ip, user, attempts)
			w.Header().Set(authHeaderStr, realmStr)
			http.Error(w, bodyStr, http.StatusUnauthorized)
			return
		}

		if user != config.User || password != config.Password {
			attempts := guard.fail(ip)
			log.Printf("basic auth failed from ip=%s user=%q attempts=%d", ip, user, attempts)
			w.Header().Set(authHeaderStr, realmStr)
			http.Error(w, bodyStr, http.StatusUnauthorized)
			return
		}

		guard.success(ip)
		h.ServeHTTP(w, r)
	})
}

type ipState struct {
	tokens         float64
	lastRefill     time.Time
	failedAttempts int
	lockUntil      time.Time
}

type ipGuard struct {
	mu      sync.Mutex
	clients map[string]*ipState
}

func newIPGuard() *ipGuard {
	return &ipGuard{
		clients: make(map[string]*ipState),
	}
}

func (g *ipGuard) get(ip string) *ipState {
	now := time.Now()
	client := g.clients[ip]
	if client == nil {
		client = &ipState{
			tokens:     rateLimitBurst,
			lastRefill: now,
		}
		g.clients[ip] = client
	}
	return client
}

func (g *ipGuard) allow(ip string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	client := g.get(ip)
	now := time.Now()
	elapsed := now.Sub(client.lastRefill).Seconds()
	client.tokens = math.Min(rateLimitBurst, client.tokens+(elapsed*rateLimitPerSecond))
	client.lastRefill = now

	if client.tokens < 1 {
		return false
	}

	client.tokens--
	return true
}

func (g *ipGuard) locked(ip string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	client := g.get(ip)
	return time.Now().Before(client.lockUntil)
}

func (g *ipGuard) fail(ip string) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	client := g.get(ip)
	client.failedAttempts++
	if client.failedAttempts >= lockoutThreshold {
		client.lockUntil = time.Now().Add(lockoutDuration)
	}
	return client.failedAttempts
}

func (g *ipGuard) success(ip string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	client := g.get(ip)
	client.failedAttempts = 0
	client.lockUntil = time.Time{}
}

func sourceIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	if host == "" {
		return "unknown"
	}
	return host
}
