package actions

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/config"
)

func SongRequest(cfg config.Settings) bot.Handler {
	requestCmd := cfg.MustString("requestCmd")

	return bot.HandlerFunc(func(s *bot.Sender) {
		var videoID string

		fields := strings.Fields(s.Msg.Text)
		for _, field := range fields {
			if strings.HasPrefix(field, "https") {
				videoID = getYoutubeVideoID(field)
				break
			}
		}

		if videoID == "" {
			log.Printf("link not found, msg: %s\n", s.Msg.Text)
			return
		}

		text := fmt.Sprintf("%s %s", requestCmd, videoID)
		s.Send(text)
	})
}

func getYoutubeVideoID(rawURL string) (videoID string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	switch u.Hostname() {
	case "www.youtube.com", "m.youtube.com":
		videoID = u.Query().Get("v")
	case "youtu.be":
		videoID = strings.Trim(u.Path, "/")
	}

	return
}
