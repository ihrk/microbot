package actions

import (
	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/config"
)

func Print(cfg config.Settings) bot.Handler {
	text := cfg.MustString("text")

	return bot.HandlerFunc(func(s *bot.Sender) {
		s.Reply(text)
	})
}
