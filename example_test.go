package oauth2_test

import (
	"io/ioutil"
	"net/http"

	"github.com/berfarah/gobot"
	"github.com/berfarah/gobot-oauth2"
	"golang.org/x/oauth2/google"
)

type ExampleChat struct{}

func (e ExampleChat) Load(*gobot.Robot)              {}
func (e ExampleChat) Username() string               { return "" }
func (e ExampleChat) Send(gobot.Message) error       { return nil }
func (e ExampleChat) Reply(gobot.Message) error      { return nil }
func (e ExampleChat) Topic(gobot.Message) error      { return nil }
func (e ExampleChat) Messages() <-chan gobot.Message { return make(chan gobot.Message) }

func Example_google() {
	oauth2.New(oauth2.Options{
		Name:         "google",
		URL:          "http://localhost:4567",
		ClientID:     "id",
		ClientSecret: "secret",
		Scopes:       []string{"scope"},
		Endpoint:     google.Endpoint,
	})
}

func ExampleAuth() {
	r := gobot.New(
		ExampleChat{},
		oauth2.New(oauth2.Options{
			Name:     "google",
			Endpoint: google.Endpoint,
			// ...
		}),
	)

	r.Hear(gobot.Contains("auth me!"), func(r gobot.Responder) error {
		var auths oauth2.Plugin
		if ok := r.Plugin(&auths); !ok {
			return nil
		}
		auths.Auth("google", r, func(c *http.Client, err error) {
			if err != nil {
				return
			}
			info, err := c.Get("https://www.googleapis.com/oauth2/v3/userinfo")
			if err != nil {
				return
			}
			defer info.Body.Close()
			b, _ := ioutil.ReadAll(info.Body)
			r.Send(gobot.Message{Text: string(b)})
		})
		return nil
	})

	r.Run()
}
