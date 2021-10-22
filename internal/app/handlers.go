package app

import (
	"log"
	"sync"
	"time"

	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/bot/actions"
	"github.com/ihrk/microbot/internal/bot/middlewares"
	"github.com/ihrk/microbot/internal/config"
	"github.com/ihrk/microbot/internal/irc"
)

func appHandler(cfg *config.App) bot.Handler {
	m := bot.NewMux(bot.MatchChannel)

	for _, ch := range cfg.Channels {
		if ch.Chat == nil {
			continue
		}

		m.Add(ch.Name, chatHandler(ch.Chat))
	}

	if cfg.Debug {
		return bot.Wrap(m, debug)
	}

	return m
}

func debug(next bot.Handler) bot.Handler {
	return bot.HandlerFunc(func(s *bot.Sender) {
		log.Println(s.Msg.Raw)
		next.Serve(s)
	})
}

func chatHandler(cfg *config.Chat) bot.Handler {
	m := bot.NewMux(bot.MatchType)

	var chatMsgs bot.Handler = bot.NewRouter(
		mux(cfg.Rewards, bot.MatchReward),
		mux(cfg.Commands, bot.MatchCmd),
	)

	chatMsgs = bot.Wrap(chatMsgs, middlewares.New(cfg.Middlewares))

	if cfg.Spam != nil {
		chatMsgs = bot.Wrap(chatMsgs, spam(cfg.Spam))
	}

	m.Add(irc.MsgTypePrivMsg, chatMsgs)

	return m
}

func mux(cfgs []*config.Trigger,
	matcher func(*irc.Msg) (string, bool),
) *bot.Mux {
	m := bot.NewMux(matcher)

	for _, cfg := range cfgs {
		m.Add(
			cfg.Key,
			actions.New(cfg.Action),
			middlewares.New(cfg.Middlewares),
		)
	}

	return m
}

type cooldown struct {
	m sync.Mutex
	t time.Time
	d time.Duration
}

func newCooldown(d time.Duration) *cooldown {
	return &cooldown{
		t: time.Now(),
		d: d,
	}
}

func (cd *cooldown) check() bool {
	var ok bool

	cd.m.Lock()

	if now := time.Now(); cd.t.Add(cd.d).Before(now) {
		cd.t = now
		ok = true
	}

	cd.m.Unlock()

	return ok
}

func spam(cfg config.Settings) bot.Middleware {
	text := cfg.MustString("text")

	cd := newCooldown(cfg.MustDuration("period"))

	return func(next bot.Handler) bot.Handler {
		return bot.HandlerFunc(func(s *bot.Sender) {
			if cd.check() {
				s.Send(text)
			}
			next.Serve(s)
		})
	}
}
