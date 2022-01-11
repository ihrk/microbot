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
	"time"

	"github.com/disintegration/gift"
	"github.com/ihrk/dots"
	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/cache"
	"github.com/ihrk/microbot/internal/config"
	"github.com/ihrk/microbot/internal/extra/bttv"
	"github.com/ihrk/microbot/internal/extra/ffz"
	"github.com/ihrk/microbot/internal/irc"
)

func Draw(_ config.Settings) bot.Handler {
	ees := newExtensionEmoteStorage()

	return bot.HandlerFunc(func(s *bot.Sender) {
		emoteURL, ok := getTwitchEmoteURL(s.Msg.Tags)
		if !ok {
			emoteURL, ok = ees.getURL(s.Msg)
		}

		if !ok {
			s.Reply("Error: emote not found")
			log.Printf("emote not found in msg: %v\n", s.Msg.Raw)
			return
		}

		var err error

		defer func() {
			if err != nil {
				s.Reply("Error: try again later")
			}
		}()

		resp, err := http.Get(emoteURL)
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

		for _, w := range strings.Fields(s.Msg.Text) {
			if w == "-evil" {
				g.Add(gift.Invert())
			}
		}

		dst := image.NewRGBA(g.Bounds(src.Bounds()))

		g.Draw(dst, src)

		p := dots.ErrDiffDithering(dst, dots.Atkinson)

		p = p.SubImage(image.Rect(0, 0, 30, 15))

		s.Reply(p.String())
	})
}

func getTwitchEmoteID(tags map[string]string) (string, bool) {
	emotes := tags["emotes"]
	end := strings.Index(emotes, ":")
	if end == -1 {
		return "", false
	}
	return emotes[:end], true
}

const urlPattern = "https://static-cdn.jtvnw.net/emoticons/v2/%v/default/dark/3.0"

func getTwitchEmoteURL(tags map[string]string) (string, bool) {
	emoteID, ok := getTwitchEmoteID(tags)
	if !ok {
		return "", false
	}

	return fmt.Sprintf(urlPattern, emoteID), true
}

type extensionEmoteStorage struct {
	c cache.Cache
}

func newExtensionEmoteStorage() *extensionEmoteStorage {
	return &extensionEmoteStorage{cache.New()}
}

func (s *extensionEmoteStorage) getURL(msg *irc.Msg) (string, bool) {
	roomID, ok := msg.Tags["room-id"]
	if !ok {
		log.Printf("room id not found in message: %s\n", msg.Raw)
		return "", false
	}

	var emoteMap map[string]emote

	v, ok := s.c.Get(roomID)
	if !ok {
		emoteMap, ok = s.collectEmotes(roomID)
		if !ok {
			return "", false
		}

		s.c.Set(roomID, emoteMap, time.Hour)
	} else {
		emoteMap = v.(map[string]emote)
	}

	words := strings.Fields(msg.Text)
	for _, word := range words {
		if e, ok := emoteMap[word]; ok {
			return e.ImageURL(), true
		}
	}

	return "", false
}

func (s *extensionEmoteStorage) collectEmotes(roomID string) (map[string]emote, bool) {
	log.Println("collecting extension emotes...")

	bttvGlobalEmotes, err := bttv.GetGlobalEmotes()
	if err != nil {
		log.Printf("bttv global emotes request failed with error: %v", err)
		return nil, false
	}

	ffzGlobalEmotes, err := ffz.GetGlobalEmotes()
	if err != nil {
		log.Printf("ffz global emotes request failed with error: %v", err)
		return nil, false
	}

	bttvUserEmotes, err := bttv.GetUserEmotes(roomID)
	if err != nil {
		log.Printf("bttv user emotes request failed with error: %v", err)
		return nil, false
	}

	ffzUserEmotes, err := ffz.GetUserEmotes(roomID)
	if err != nil {
		log.Printf("ffz user emotes request failed with error: %v", err)
		return nil, false
	}

	emoteMap := map[string]emote{}

	for _, e := range bttvGlobalEmotes {
		emoteMap[e.Code] = e
	}

	for _, e := range ffzGlobalEmotes {
		emoteMap[e.Code] = e
	}

	for _, e := range bttvUserEmotes.ChannelEmotes {
		emoteMap[e.Code] = e
	}

	for _, e := range bttvUserEmotes.SharedEmotes {
		emoteMap[e.Code] = e
	}

	for _, e := range ffzUserEmotes {
		emoteMap[e.Code] = e
	}

	return emoteMap, true
}

type emote interface {
	ImageURL() string
}
