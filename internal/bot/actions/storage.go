package actions

import (
	"log"

	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/config"
)

type Storage map[string]func(cfg config.Settings) bot.Handler

var defaultStorage = Storage{
	"print":       Print,
	"elo":         Elo,
	"songRequest": SongRequest,
	"draw":        Draw,
}

func New(cfg *config.Feature) bot.Handler {
	b, ok := defaultStorage[cfg.Type]
	if !ok {
		log.Fatalf("unknown action: %v\n", cfg.Type)
	}

	return b(cfg.Settings)
}
