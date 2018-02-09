package oauth2_test

import (
	"io/ioutil"
	"net/http"

	"github.com/botopolis/bot"
	"github.com/botopolis/bot/mock"
	"github.com/botopolis/oauth2"
	"golang.org/x/oauth2/google"
)

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

func Example() {
	ExampleChat := mock.NewChat()
	ExampleChat.MessageChan = make(chan bot.Message)
	close(ExampleChat.MessageChan)

	r := bot.New(
		ExampleChat,
		oauth2.New(oauth2.Options{
			Name:     "google",
			URL:      "http://localhost:4567",
			Endpoint: google.Endpoint,
			// ...
		}),
	)

	r.Hear(bot.Contains("auth me!"), func(r bot.Responder) error {
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
			r.Send(bot.Message{Text: string(b)})
		})
		return nil
	})

	r.Run()
	// Output:
}
