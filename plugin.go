package oauth2

import (
	"fmt"
	"net/http"

	"github.com/botopolis/bot"
)

// Plugin is the oauth2 plugin
type Plugin struct {
	opts       []Options
	strategies map[string]*Strategy
}

// New creates a new oauth2 plugin with strategies
// corresponding to options passed in
func New(o ...Options) *Plugin {
	return &Plugin{
		opts:       o,
		strategies: make(map[string]*Strategy),
	}
}

// Load is called when the plugin is loaded by gobto.Robot
// This is the step at which options are validated
func (p *Plugin) Load(r *bot.Robot) {
	for _, o := range p.opts {
		if _, ok := p.strategies[o.Name]; ok {
			r.Logger.Errorf("Duplicate oauth2 strategies registered: %s", o.Name)
			continue
		}
		s := &Strategy{Opts: o}
		s.Load(r)
		p.strategies[o.Name] = s
	}
}

// Auth delegates to Strategy.Auth()
func (p *Plugin) Auth(strategy string, r bot.Responder, f func(*http.Client, error)) {
	if s, ok := p.strategies[strategy]; ok {
		s.Auth(r, f)
		return
	}
	f(nil, fmt.Errorf("Unregistered strategy: %s", strategy))
}
