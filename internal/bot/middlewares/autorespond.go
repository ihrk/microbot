package middlewares

import (
	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/config"
	"github.com/ihrk/microbot/internal/cooldown"
)

func Autorespond(cfg config.Settings) bot.Middleware {
	text := cfg.MustString("text")

	per := cfg.MustDuration("period")
	gap, _ := cfg.Int("gap")

	cd := cooldown.New(per, gap)

	return func(next bot.Handler) bot.Handler {
		return bot.HandlerFunc(func(s *bot.Sender) {
			if cd.Check() {
				s.Send(text)
			}
			next.Serve(s)
		})
	}
}
