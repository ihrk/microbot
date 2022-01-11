package bot

import (
	"strings"

	"github.com/ihrk/microbot/internal/irc"
)

type HandlerFunc func(s *Sender)

var _ Handler = (HandlerFunc)(nil)

func (f HandlerFunc) Serve(s *Sender) {
	f(s)
}

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

type Router interface {
	Match(*irc.Msg) (Handler, bool)
}

type StringRouter struct {
	middleware Middleware
	matcher    func(*irc.Msg) (string, bool)
	handlers   map[string]Handler
}

var _ Router = (*StringRouter)(nil)

func NewStringRouter(
	matcher func(*irc.Msg) (string, bool),
	middlewares ...Middleware,
) *StringRouter {
	return &StringRouter{
		middleware: Concat(middlewares...),
		matcher:    matcher,
		handlers:   make(map[string]Handler),
	}
}

func (r *StringRouter) Match(msg *irc.Msg) (Handler, bool) {
	key, ok := r.matcher(msg)
	if !ok {
		return nil, false
	}

	h, ok := r.handlers[key]

	return h, ok
}

func (r *StringRouter) Add(
	key string,
	handler Handler,
	middlewares ...Middleware,
) {
	r.handlers[key] = Wrap(
		handler,
		r.middleware,
		Concat(middlewares...),
	)
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

type SingleRouter struct {
	h Handler
}

var _ Router = (*SingleRouter)(nil)

func NewSingleRouter(h Handler) *SingleRouter {
	return &SingleRouter{h}
}

func (r *SingleRouter) Match(_ *irc.Msg) (Handler, bool) {
	return r.h, true
}

type Mux struct {
	routers []Router
}

var _ Handler = (*Mux)(nil)

func NewMux(routers ...Router) Handler {
	return &Mux{
		routers: routers,
	}
}

func (m *Mux) Serve(s *Sender) {
	for _, r := range m.routers {
		h, ok := r.Match(s.Msg)
		if ok {
			h.Serve(s)
			return
		}
	}
}
