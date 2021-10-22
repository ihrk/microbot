package bot

import (
	"strings"

	"github.com/ihrk/microbot/internal/irc"
)

type HandlerFunc func(s *Sender)

func (f HandlerFunc) Serve(s *Sender) {
	f(s)
}

var _ Handler = (HandlerFunc)(nil)

func Concat(middlewares ...Middleware) Middleware {
	return func(h Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
}

func Wrap(handler Handler, middlewares ...Middleware) Handler {
	return Concat(middlewares...)(handler)
}

var _ Handler = (*Mux)(nil)

type Mux struct {
	middleware Middleware
	matcher    func(*irc.Msg) (string, bool)
	handlers   map[string]Handler
}

func NewMux(matcher func(*irc.Msg) (string, bool), middlewares ...Middleware) *Mux {
	return &Mux{
		middleware: Concat(middlewares...),
		matcher:    matcher,
		handlers:   make(map[string]Handler),
	}
}

func (m *Mux) Add(key string, handler Handler, middlewares ...Middleware) {
	m.handlers[key] = Wrap(handler, m.middleware, Concat(middlewares...))
}

func (m *Mux) Serve(s *Sender) {
	if key, ok := m.matcher(s.Msg); ok {
		if handler, ok := m.handlers[key]; ok {
			handler.Serve(s)
		}
	}
}

func MatchType(msg *irc.Msg) (string, bool) {
	return msg.Type, msg.Type != ""
}

func MatchChannel(msg *irc.Msg) (string, bool) {
	return msg.Channel, msg.Channel != ""
}

func MatchCmd(msg *irc.Msg) (string, bool) {
	if len(msg.Text) < 2 || msg.Text[0] != '!' {
		return "", false
	}

	end := strings.Index(msg.Text, " ")
	if end == -1 {
		end = len(msg.Text)
	}

	return msg.Text[1:end], true
}

func MatchReward(msg *irc.Msg) (string, bool) {
	rewardID, ok := msg.Tags["custom-reward-id"]
	return rewardID, ok
}

var _ Handler = (*Router)(nil)

type Router struct {
	muxes []*Mux
}

func NewRouter(muxes ...*Mux) *Router {
	return &Router{
		muxes: muxes,
	}
}

func (r *Router) Serve(s *Sender) {
	for _, m := range r.muxes {
		if key, ok := m.matcher(s.Msg); ok {
			if handler, ok := m.handlers[key]; ok {
				handler.Serve(s)
				return
			}
		}
	}
}
