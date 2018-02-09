package oauth2

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/botopolis/bot"
	"github.com/botopolis/bot/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestAuth_cacheHit(t *testing.T) {
	var run bool
	s := testStrategy("")
	store := newStore("", bot.NewBrain())
	store.Set("Jean", oauth2.Token{})
	s.store = store
	res := bot.Responder{}
	res.User = "Jean"

	s.Auth(res, func(c *http.Client, err error) {
		run = true
		assert.NotNil(t, c)
		assert.Nil(t, err)
	})
	assert.True(t, run)
}

func TestAuth_cacheMiss(t *testing.T) {
	var run bool
	s := testStrategy("")
	s.store = newStore("", bot.NewBrain())
	res := bot.Responder{Robot: bot.New(mock.NewChat())}

	cb := func(c *http.Client, err error) { run = true }
	s.Auth(res, cb)

	assert.False(t, run, "callback is not evaluated")

	for _, as := range s.authSessions.registry {
		as.Run(nil, nil)
	}

	assert.True(t, run, "callback is saved in authSessions")
}

func TestHandleLogin(t *testing.T) {
	s := testStrategy("")
	cases := []struct {
		ReqURL  string
		ResCode int
		ResURL  string
	}{
		{
			ReqURL:  "http://foo.com/login?state=foo",
			ResCode: http.StatusFound,
			ResURL:  "/auth?client_id=&redirect_uri=http%3A%2F%2Ffoo.com%2Fauth&response_type=code&state=foo",
		},
		{
			ReqURL:  "http://foo.com/login",
			ResCode: http.StatusBadRequest,
			ResURL:  "",
		},
	}

	for _, c := range cases {
		requestURL, err := url.Parse(c.ReqURL)
		recorder := httptest.NewRecorder()

		s.HandleLogin(recorder, &http.Request{URL: requestURL})

		assert.Nil(t, err)
		assert.Equal(t, c.ResCode, recorder.Code)
		assert.Equal(t, c.ResURL, recorder.Header().Get("Location"))
	}
}

func TestHandleAuth_success(t *testing.T) {
	d := handleAuthTestData{State: "foo", URL: "http://foo.com?state=foo&code=secret"}
	ts := d.Server(http.StatusOK)

	s := testStrategy(ts.URL)
	s.store = newStore("", bot.NewBrain())
	s.authSessions.Set(d.State, authSession{
		Func: func(c *http.Client, err error) {
			d.Ran = true
			assert.Nil(t, err)
		},
		User: "Jean",
	})

	// under test
	s.HandleAuth(d.Recorder(), d.Request())
	assert.True(t, d.Ran)
	assert.Equal(t, http.StatusOK, d.Recorder().Code)
}

func TestHandleAuth_invalidState(t *testing.T) {
	d := handleAuthTestData{State: "foo", URL: "http://foo.com?state=foo&code=secret"}
	ts := d.Server(http.StatusOK)
	defer ts.Close()
	s := testStrategy(ts.URL)

	// under test
	s.HandleAuth(d.Recorder(), d.Request())
	assert.Equal(t, http.StatusUnauthorized, d.Recorder().Code, "Returns Unauthorized code when state isn't in our session cache")
}

func TestHandleAuth_failedExchange(t *testing.T) {
	// setup
	d := handleAuthTestData{State: "foo", URL: "http://foo.com?state=foo&code=secret"}
	ts := d.Server(http.StatusBadGateway)
	defer ts.Close()
	s := testStrategy(ts.URL)
	s.authSessions.Set(d.State, authSession{Func: func(c *http.Client, err error) {
		d.Ran = true
		assert.Nil(t, c)
		assert.NotNil(t, err)
	}})

	// under test
	s.HandleAuth(d.Recorder(), d.Request())
	assert.True(t, d.Ran)
	assert.Equal(t, http.StatusBadRequest, d.Recorder().Code, "Returns Unauthorized code when state isn't in our session cache")
}

type handleAuthTestData struct {
	Ran      bool
	State    string
	URL      string
	recorder *httptest.ResponseRecorder
}

func (d *handleAuthTestData) Request() *http.Request {
	url, _ := url.Parse(d.URL)
	return &http.Request{URL: url}
}

func (d *handleAuthTestData) Recorder() *httptest.ResponseRecorder {
	if d.recorder == nil {
		d.recorder = httptest.NewRecorder()
	}
	return d.recorder
}

func (d *handleAuthTestData) Server(statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if statusCode == 200 {
			w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
			w.Write([]byte("access_token=90d64460d14870c08c81352a05dedd3465940a7c&scope=user&token_type=bearer"))
			return
		}
		w.WriteHeader(statusCode)
	}))
}

func testStrategy(url string) Strategy {
	return Strategy{
		Opts: Options{},
		Config: &oauth2.Config{
			RedirectURL: "http://foo.com/auth",
			Endpoint: oauth2.Endpoint{
				AuthURL:  url + "/auth",
				TokenURL: url + "/token",
			},
		},
	}
}
