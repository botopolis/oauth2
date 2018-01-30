package oauth2_test

import (
	"testing"

	"github.com/botopolis/oauth2"
	"github.com/stretchr/testify/assert"
)

var validateTestcases = []struct {
	Name  string
	URL   string
	Error string
}{
	{
		Name:  "",
		URL:   "",
		Error: "Name is a required",
	},
	{
		Name:  "foo",
		URL:   "",
		Error: "URL is a required",
	},
	{
		Name:  "foo;",
		URL:   "foo.com",
		Error: "Name may not contain",
	},
	{
		Name:  "foo/",
		URL:   "foo.com",
		Error: "Name may not contain",
	},
	{
		Name:  "foo",
		URL:   "foo.com",
		Error: "URL must be valid",
	},
	{
		Name:  "foo",
		URL:   "http://foo.com",
		Error: "",
	},
}

func TestOptions_Validate(t *testing.T) {
	for _, c := range validateTestcases {
		o := oauth2.Options{Name: c.Name, URL: c.URL}
		err := o.Validate()
		if c.Error == "" {
			assert.Nil(t, err)
			continue
		}
		assert.Contains(t, err.Error(), c.Error)
	}
}

func TestOptions_AuthPath(t *testing.T) {
	name := "foo"
	o := oauth2.Options{Name: name}
	assert.Equal(t, "/oauth2/foo/auth", o.AuthPath())
}

func TestOptions_AuthURL(t *testing.T) {
	name := "foo"
	url := "http://bar.com"
	o := oauth2.Options{Name: name, URL: url}
	assert.Equal(t, "http://bar.com/oauth2/foo/auth", o.AuthURL())
}

func TestOptions_LoginPath(t *testing.T) {
	name := "foo"
	o := oauth2.Options{Name: name}
	assert.Equal(t, "/oauth2/foo/login", o.LoginPath())
}

func TestOptions_LoginURL(t *testing.T) {
	name := "foo"
	url := "http://bar.com"
	o := oauth2.Options{Name: name, URL: url}
	assert.Equal(t, "http://bar.com/oauth2/foo/login", o.LoginURL())
}
