package actions

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strings"

	"github.com/disintegration/gift"
	"github.com/ihrk/dots"
	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/config"
)

func Draw(_ config.Settings) bot.Handler {
	return bot.HandlerFunc(func(s *bot.Sender) {
		emoteID, ok := getEmoteID(s.Msg.Tags)
		if !ok {
			log.Printf("emote not found in msg: %v\n", s.Msg.Raw)
			return
		}

		var err error

		defer func() {
			if err != nil {
				s.Reply("Error: try again later")
			}
		}()

		resp, err := http.Get(getEmoteURL(emoteID))
		if err != nil {
			log.Printf("get emote request error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		var src image.Image
		switch tp := resp.Header.Get("Content-Type"); tp {
		case "image/png":
			src, err = png.Decode(resp.Body)
		case "image/gif":
			src, err = gif.Decode(resp.Body)
		case "image/jpeg":
			src, err = jpeg.Decode(resp.Body)
		default:
			err = image.ErrFormat
		}

		if err != nil {
			log.Printf("image decode error: %v\n", err)
			return
		}

		g := gift.New(
			gift.UnsharpMask(4, 2, 0),
			gift.Resize(60, 0, gift.LanczosResampling),
			gift.UnsharpMask(2, 1, 0),
		)

		dst := image.NewRGBA(g.Bounds(src.Bounds()))

		g.Draw(dst, src)

		p := dots.ErrDiffDithering(dst, dots.Atkinson)

		p = p.SubImage(image.Rect(0, 0, 30, 15))

		s.Reply(p.String())
	})
}

func getEmoteID(tags map[string]string) (string, bool) {
	emotes := tags["emotes"]
	end := strings.Index(emotes, ":")
	if end == -1 {
		return "", false
	}
	return emotes[:end], true
}

const urlPattern = "https://static-cdn.jtvnw.net/emoticons/v2/%v/default/dark/3.0"

func getEmoteURL(emoteID string) string {
	return fmt.Sprintf(urlPattern, emoteID)
}
