package actions

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/config"
)

var urlExpr = regexp.MustCompile(`(http(s)?:\/\/.)?(www\.)?[-a-zA-Z0-9@:%._+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_+.~#?&/=]*)`)

func SongRequest(cfg config.Settings) bot.Handler {
	requestCmd := cfg.MustString("requestCmd")

	return bot.HandlerFunc(func(s *bot.Sender) {
		rawURL := urlExpr.FindString(s.Msg.Text)
		videoID, err := getYoutubeVideoID(rawURL)
		if err != nil {
			log.Printf("link not found with err: %s, msg: %s\n", err, s.Msg.Text)
			return
		}

		text := fmt.Sprintf("%s %s", requestCmd, videoID)
		s.Send(text)
	})
}

func getYoutubeVideoID(rawURL string) (string, error) {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", err
	}

	var videoID string

	switch host := u.Hostname(); host {
	case "www.youtube.com", "m.youtube.com":
		videoID = u.Query().Get("v")
	case "youtu.be":
		videoID = strings.Trim(u.Path, "/")
	default:
		err = fmt.Errorf("unknown hostname: %s", host)
	}

	return videoID, err
}
