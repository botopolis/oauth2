package oauth2

import (
	"testing"

	"github.com/botopolis/bot"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestStore(t *testing.T) {
	s := newStore("test", bot.NewBrain())
	token := oauth2.Token{AccessToken: "success"}

	_, ok := s.Get("foo")
	assert.False(t, ok, "doesn't find anything in an empty store")

	s.Set("bob", token)
	out, ok := s.Get("bob")
	assert.Equal(t, token, out, "matches what we put in")
	assert.True(t, ok, "finds the value")

	s.Delete("bob")
	_, ok = s.Get("bob")
	assert.False(t, ok, "doesn't find deleted items")
}
