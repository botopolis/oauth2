package oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/botopolis/bot"
	"golang.org/x/oauth2"
)

// Strategy is an OAuth2 strategy
type Strategy struct {
	// Options passed in
	Opts Options
	// Configuration for strategy
	Config *oauth2.Config

	// config holds the ouath2 config
	store        *store
	authSessions authSessions
}

// Load initializes the strategy and conforms to bot.Plugin interface
func (s *Strategy) Load(r *bot.Robot) {
	o := s.Opts
	if err := o.Validate(); err != nil {
		r.Logger.Errorf("%s: %s", o.Name, err.Error())
	}
	s.store = newStore(o.Name, r.Brain)
	s.Config = &oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		Scopes:       o.Scopes,
		Endpoint:     o.Endpoint,
		RedirectURL:  o.AuthURL(),
	}

	r.Router.HandleFunc(o.AuthPath(), s.HandleAuth).Methods("GET")
	r.Router.HandleFunc(o.LoginPath(), s.HandleLogin).Methods("GET")
}

// Auth is meant to be called by other plugins
func (s *Strategy) Auth(r bot.Responder, f func(*http.Client, error)) {
	if token, ok := s.store.Get(r.User); ok {
		f(s.Config.Client(context.Background(), &token), nil)
		return
	}

	sid := randToken()
	err := r.Direct(fmt.Sprintf(
		"I need you to log in before you can do that: %s?state=%s",
		s.Opts.LoginURL(),
		sid,
	))

	if err != nil {
		f(nil, err)
		return
	}

	s.authSessions.Set(sid, authSession{Func: f, User: r.User})
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// HandleAuth handles the oauth2 callback and finds the user by
// looking up the state in the cache (authSessions)
func (s *Strategy) HandleAuth(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	state := query.Get("state")
	code := query.Get("code")

	session, ok := s.authSessions.Get(state)
	if !ok {
		session.Run(nil, errors.New("Invalid state in callback"))
		http.Error(w, "Invalid state", http.StatusUnauthorized)
		return
	}
	s.authSessions.Delete(state)

	token, err := s.Config.Exchange(context.Background(), code)
	if err != nil {
		session.Run(nil, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.store.Set(session.User, *token)
	session.Run(s.Config.Client(context.Background(), token), nil)
}

// HandleLogin redirects the user to the AuthCode URL
func (s *Strategy) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, s.Config.AuthCodeURL(state), 302)
}
