package app

import (
	"log"

	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/bot/actions"
	"github.com/ihrk/microbot/internal/bot/middlewares"
	"github.com/ihrk/microbot/internal/config"
	"github.com/ihrk/microbot/internal/irc"
)

func appHandler(cfg *config.App) bot.Handler {
	r := bot.NewStringRouter(bot.MatchChannel)

	for _, ch := range cfg.Channels {
		if ch.Chat == nil {
			continue
		}

		r.Add(ch.Name, chatHandler(ch.Chat))
	}

	m := bot.NewMux(r)

	if cfg.Debug {
		m = bot.Wrap(m, debug)
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
	chat := bot.NewMux(
		newRouter(cfg.Rewards, bot.MatchReward),
		newRouter(cfg.Commands, bot.MatchCmd),
	)

	r := bot.NewStringRouter(bot.MatchType)
	r.Add(
		irc.MsgTypePrivMsg,
		chat,
		middlewares.New(cfg.Middlewares),
	)

	return bot.NewMux(r)
}

func newRouter(
	cfgs []*config.Trigger,
	matcher func(*irc.Msg) (string, bool),
) bot.Router {
	r := bot.NewStringRouter(matcher)

	for _, cfg := range cfgs {
		r.Add(
			cfg.Key,
			actions.New(cfg.Action),
			middlewares.New(cfg.Middlewares),
		)
	}

	return r
}
