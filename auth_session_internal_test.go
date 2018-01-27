package oauth2

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthSession(t *testing.T) {
	s := newAuthSessions()

	session := authSession{User: "bob", Func: func(*http.Client, error) {}}
	sid := "foo"

	_, ok := s.Get("bar")
	assert.False(t, ok, "doesn't find anything in an empty store")

	s.Set(sid, session)
	out, ok := s.Get(sid)
	assert.Equal(t, session.User, out.User, "matches what we put in")
	assert.True(t, ok, "finds the value")

	s.Delete(sid)
	_, ok = s.Get(sid)
	assert.False(t, ok, "doesn't find deleted items")
}
