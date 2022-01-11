package middlewares

import (
	"log"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/config"
	"github.com/ihrk/microbot/internal/irc"
	"github.com/ihrk/microbot/internal/limit"
)

type filter struct {
	ff filterFunc
	h  bot.Handler

	passThrough bool

	allowMod     bool
	allowVIP     bool
	allowSub     bool
	allowRewards []string
}

type filterFunc func(*irc.Msg) bool

func Filter(cfg config.Settings) bot.Middleware {
	var f filter

	f.h = newFilterHandler(cfg)
	f.ff = newFilterFunc(cfg)

	f.passThrough = cfg.Bool("passThrough")

	f.allowMod = cfg.Bool("allowMod")
	f.allowVIP = cfg.Bool("allowVIP")
	f.allowSub = cfg.Bool("allowSub")

	f.allowRewards = cfg.Strings("allowRewards")

	return f.mw
}

func (f *filter) mw(next bot.Handler) bot.Handler {
	return bot.HandlerFunc(func(s *bot.Sender) {
		rewardUUID, isReward := getReward(s.Msg)

		ok := isBroadcaster(s.Msg) ||
			f.allowMod && isMod(s.Msg) ||
			f.allowVIP && isVIP(s.Msg) ||
			f.allowSub && isSub(s.Msg) ||
			isReward && elem(rewardUUID, f.allowRewards) ||
			f.ff(s.Msg)

		if !ok {
			f.h.Serve(s)
		}

		if ok || f.passThrough {
			next.Serve(s)
		}
	})
}

func pureFunc(ff filterFunc) func(config.Settings) filterFunc {
	return func(_ config.Settings) filterFunc {
		return ff
	}
}

var filterFuncStorage = map[string]func(config.Settings) filterFunc{
	"byUsername": byUsername,
	"countLimit": countLimit,
	"limitChars": limitChars,
	"blockLinks": pureFunc(blockLinks),
	"blockAll":   pureFunc(blockAll),
}

func newFilterFunc(cfg config.Settings) filterFunc {
	filterType := cfg.MustString("type")

	bf, found := filterFuncStorage[filterType]
	if !found {
		log.Fatalf("filter type not found: %s\n", filterType)
	}

	return bf(cfg)
}

const (
	penaltyDeleteMsg = "deleteMsg"
	penaltyTimeout   = "timeout"
	penaltyBan       = "ban"
)

var penaltyTypes = []string{
	penaltyDeleteMsg,
	penaltyTimeout,
	penaltyBan,
}

func newFilterHandler(cfg config.Settings) bot.Handler {
	penalty, _ := cfg.StringFromSet("penalty", penaltyTypes)

	var duration time.Duration
	if penalty == penaltyTimeout {
		duration = cfg.MustDuration("duration")
	}

	replyText, hasReply := cfg.String("reply")

	reason, _ := cfg.String("reason")

	return bot.HandlerFunc(func(s *bot.Sender) {
		switch penalty {
		case penaltyDeleteMsg:
			s.Delete()
		case penaltyTimeout:
			s.Timeout(duration, reason)
		case penaltyBan:
			s.Ban(reason)
		}

		if hasReply {
			s.Reply(replyText)
		}
	})
}

func countLimit(cfg config.Settings) filterFunc {
	lim := cfg.MustInt("limitAmount")
	per := cfg.MustDuration("limitPeriod")

	l := limit.New(lim, per)
	return func(msg *irc.Msg) bool {
		return l.Add(1)
	}
}

func isBroadcaster(msg *irc.Msg) bool {
	return strings.Contains(msg.Tags["badges"], "broadcaster")
}

func isMod(msg *irc.Msg) bool {
	return strings.Contains(msg.Tags["badges"], "moderator")
}

func isVIP(msg *irc.Msg) bool {
	return strings.Contains(msg.Tags["badges"], "vip")
}

func isSub(msg *irc.Msg) bool {
	return strings.Contains(msg.Tags["badges"], "subscriber")
}

func getReward(msg *irc.Msg) (string, bool) {
	s, ok := msg.Tags["custom-reward-id"]
	return s, ok
}

func elem(s string, a []string) bool {
	for i := range a {
		if s == a[i] {
			return true
		}
	}

	return false
}

func byUsername(cfg config.Settings) filterFunc {
	name := cfg.MustString("username")

	return func(msg *irc.Msg) bool {
		return msg.User == name
	}
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

func limitChars(cfg config.Settings) filterFunc {
	tp := cfg.MustString("charType")
	var charFunc func(rune) bool
	switch tp {
	case "symbol":
		charFunc = unicode.IsSymbol
	case "upper":
		charFunc = unicode.IsUpper
	case "mark":
		charFunc = unicode.IsMark
	default:
		log.Fatalf("unknown character type: %s\n", tp)
	}

	limit := cfg.MustInt("charLimit")

	return func(msg *irc.Msg) bool {
		return countFunc(msg.Text, charFunc) <= limit
	}
}

var urlExpr = regexp.MustCompile(`(http(s)?:\/\/.)?(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_+.~#?&/=]*)`)

func blockLinks(msg *irc.Msg) bool {
	return !urlExpr.MatchString(msg.Text)
}

func blockAll(_ *irc.Msg) bool {
	return false
}
