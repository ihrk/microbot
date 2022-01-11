package middlewares

import (
	"log"

	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/config"
)

type Storage map[string]func(cfg config.Settings) bot.Middleware

var defaultStorage = Storage{
	"filter":      Filter,
	"autorespond": Autorespond,
}

func New(cfgs []*config.Feature) bot.Middleware {
	mws := make([]bot.Middleware, len(cfgs))

	for i, cfg := range cfgs {
		b, ok := defaultStorage[cfg.Type]
		if !ok {
			log.Fatalf("unknown middleware: %v\n", cfg.Type)
		}
		mws[i] = b(cfg.Settings)
	}

	return bot.Concat(mws...)
}
