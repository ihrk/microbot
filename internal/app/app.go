package app

import (
	"context"
	"log"
	"time"

	"github.com/ihrk/microbot/internal/backoff"
	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/config"
	"github.com/ihrk/microbot/internal/creds"
	"github.com/ihrk/microbot/internal/irc"
)

const (
	dialTimeout    = 10 * time.Second
	retryLim       = 10
	initialBackoff = 2 * time.Second
)

func LoadConfigAndRun(configPath, credsPath string) error {
	cfg, err := config.Read(configPath)
	if err != nil {
		return err
	}

	err = creds.Load(credsPath)
	if err != nil {
		return err
	}

	return loadConfig(cfg).run(context.Background())
}

type app struct {
	h        bot.Handler
	channels []string
}

func loadConfig(cfg *config.App) *app {
	var a app

	for _, ch := range cfg.Channels {
		a.channels = append(a.channels, ch.Name)
	}

	a.h = appHandler(cfg)

	return &a
}

func (a *app) run(ctx context.Context) error {
	var (
		client *irc.Client
		err    error
	)

	for {
		err = backoff.RunWithRetry(retryLim, initialBackoff,
			func() error {
				var dialErr error
				log.Println("attempting to dial...")
				client, dialErr = irc.Dial(ctx, dialTimeout)
				return dialErr
			})
		if err != nil {
			return err
		}

		log.Println("dial is successful")

		err = a.listenAndServe(ctx, client)
		log.Printf("connection interrupted with error: %v\n", err)
	}
}

func (a *app) listenAndServe(ctx context.Context, c *irc.Client) error {
	defer c.Disconnect()

	err := c.RegCaps(irc.CapTags, irc.CapCommands)
	if err != nil {
		return err
	}

	err = c.Login(creds.TwitchUser(), creds.TwitchPass())
	if err != nil {
		return err
	}

	for _, channel := range a.channels {
		err = c.Join(channel)
		if err != nil {
			return err
		}
	}

	return bot.NewServer(c, a.h).ListenAndServe(ctx)
}
