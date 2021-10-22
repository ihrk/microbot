package middlewares

import (
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/config"
	"github.com/ihrk/microbot/internal/irc"
	"github.com/ihrk/microbot/internal/queue"
)

type filter struct {
	ff func(msg *irc.Msg) bool
	h  bot.Handler
}

func (f *filter) mw(next bot.Handler) bot.Handler {
	return bot.HandlerFunc(func(s *bot.Sender) {
		if f.ff(s.Msg) {
			next.Serve(s)
		} else {
			f.h.Serve(s)
		}
	})
}

func filterHandler(cfg config.Settings) bot.Handler {
	filterAction, _ := cfg.String("filterAction")

	var d time.Duration
	if filterAction == "timeout" {
		d = cfg.MustDuration("duration")
	}

	var replyText string
	if filterAction == "reply" {
		replyText = cfg.MustString("text")
	}

	reason, _ := cfg.String("reason")

	switch filterAction {
	case "":
		return bot.HandlerFunc(func(s *bot.Sender) {})
	case "ban":
		return bot.HandlerFunc(func(s *bot.Sender) {
			s.Ban(reason)
		})
	case "timeout":
		return bot.HandlerFunc(func(s *bot.Sender) {
			s.Timeout(d, reason)
		})
	case "reply":
		return bot.HandlerFunc(func(s *bot.Sender) {
			s.Reply(replyText)
		})
	default:
		log.Fatalf("unknown filter action: %s\n", filterAction)
	}

	return nil
}

func newFilter(cfg config.Settings, ff func(*irc.Msg) bool) bot.Middleware {
	var f filter

	f.ff = ff
	f.h = filterHandler(cfg)

	return f.mw
}

type rateLimit struct {
	lim int
	q   queue.Queue
}

func newRateLimit(lim int, d time.Duration) *rateLimit {
	return &rateLimit{
		q: queue.NewQueue(d),
	}
}

func (rl *rateLimit) ff(_ *irc.Msg) bool {
	_, ok := rl.q.TryPush(1, rl.lim)
	return ok
}

func FilterRatelimit(cfg config.Settings) bot.Middleware {
	lim := cfg.MustInt("limit")
	period := cfg.MustDuration("period")

	rl := newRateLimit(lim, period)

	return newFilter(cfg, rl.ff)
}

func isBroadcaster(msg *irc.Msg) bool {
	return strings.Contains(msg.Tags["badges"], "broadcaster")

}

func isMod(msg *irc.Msg) bool {
	return isBroadcaster(msg) ||
		strings.Contains(msg.Tags["badges"], "moderator")
}

func isVIP(msg *irc.Msg) bool {
	return isMod(msg) ||
		strings.Contains(msg.Tags["badges"], "vip")
}

func isSub(msg *irc.Msg) bool {
	return isMod(msg) ||
		strings.Contains(msg.Tags["badges"], "subscriber")
}

func FilterBroadcaster(cfg config.Settings) bot.Middleware {
	return newFilter(cfg, isBroadcaster)
}

func FilterMod(cfg config.Settings) bot.Middleware {
	return newFilter(cfg, isMod)
}

func FilterVIP(cfg config.Settings) bot.Middleware {
	return newFilter(cfg, isVIP)
}

func FilterSub(cfg config.Settings) bot.Middleware {
	return newFilter(cfg, isSub)
}

func FilterUser(cfg config.Settings) bot.Middleware {
	name := cfg.MustString("username")

	return newFilter(cfg, func(msg *irc.Msg) bool {
		return msg.User == name
	})
}

func countFunc(s string, f func(rune) bool) int {
	var n int
	for _, r := range s {
		if f(r) {
			n++
		}
	}
	return n
}

func FilterLimitChars(cfg config.Settings) bot.Middleware {
	tp := cfg.MustString("charType")
	var charFunc func(rune) bool
	switch tp {
	case "symbol":
		charFunc = unicode.IsSymbol
	case "upper":
		charFunc = unicode.IsUpper
	default:
		log.Fatalf("unknown character type: %s\n", tp)
	}

	limit := cfg.MustInt("charLimit")

	return newFilter(cfg, func(msg *irc.Msg) bool {
		return countFunc(msg.Text, charFunc) <= limit
	})
}
