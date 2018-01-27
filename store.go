package oauth2

import (
	"fmt"

	"github.com/berfarah/gobot"
	"golang.org/x/oauth2"
)

// store is responsible for storing tokens for individual
// users by oauth service
type store struct {
	namespace string
	store     *gobot.Brain
}

func newStore(namespace string, b *gobot.Brain) *store {
	return &store{
		namespace: namespace,
		store:     b,
	}
}

func (s store) key(user string) string {
	return fmt.Sprintf("auth:%s:%s", s.namespace, user)
}

func (s store) Get(user string) (oauth2.Token, bool) {
	var t oauth2.Token
	err := s.store.Get(s.key(user), &t)
	return t, err == nil
}

func (s store) Set(user string, t oauth2.Token) {
	s.store.Set(s.key(user), &t)
}

func (s store) Delete(user string) {
	s.store.Delete(s.key(user))
}
