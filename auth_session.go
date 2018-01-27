package oauth2

import (
	"net/http"
	"sync"
)

// authSessions caches auth sessions by state and
// ensures callbacks are called when authentication
// completes
type authSessions struct {
	mu       sync.Mutex
	registry map[string]authSession
}

type authSession struct {
	User string
	Func func(*http.Client, error)
}

func newAuthSessions() *authSessions {
	return &authSessions{
		registry: make(map[string]authSession),
	}
}

func (a authSession) Run(h *http.Client, err error) {
	if a.Func != nil {
		a.Func(h, err)
	}
}

func (as *authSessions) Get(key string) (authSession, bool) {
	as.mu.Lock()
	defer as.mu.Unlock()
	a, ok := as.registry[key]
	return a, ok
}

func (as *authSessions) Set(key string, a authSession) {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.registry[key] = a
}

func (as *authSessions) Delete(key string) {
	as.mu.Lock()
	defer as.mu.Unlock()
	delete(as.registry, key)
}
