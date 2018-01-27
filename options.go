package oauth2

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"golang.org/x/oauth2"
)

// Options for Oauth2 strategy
type Options struct {
	// Unique name for strategy
	Name string
	// URL for server
	URL string

	// Oauth2 options
	ClientID     string
	ClientSecret string
	Scopes       []string
	Endpoint     oauth2.Endpoint
}

func required(key, value string) error {
	if value == "" {
		return fmt.Errorf("%s is a required option for oauth2.Strategy", key)
	}
	return nil
}

const reservedChars = ";,/"

// Validate checks if the options are valid
func (o Options) Validate() error {
	if err := required("Name", o.Name); err != nil {
		return err
	}
	if err := required("URL", o.URL); err != nil {
		return err
	}

	if strings.ContainsAny(o.Name, reservedChars) {
		return fmt.Errorf("Name may not contain characters: %s", reservedChars)
	}

	u, err := url.Parse(o.URL)
	if err != nil || u.Hostname() == "" {
		return fmt.Errorf("URL must be valid")
	}
	return nil
}

// AuthPath returns the auth path
func (o Options) AuthPath() string {
	return path.Join("/oauth2", o.Name, "auth")
}

// AuthURL returns the auth URL
func (o Options) AuthURL() string {
	return o.URL + o.AuthPath()
}

// LoginPath returns the login path
func (o Options) LoginPath() string {
	return path.Join("/oauth2", o.Name, "login")
}

// LoginURL returns the login URL
func (o Options) LoginURL() string {
	return o.URL + o.LoginPath()
}
