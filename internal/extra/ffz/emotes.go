package ffz

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	globalEmotesURL = "https://api.betterttv.net/3/cached/frankerfacez/emotes/global"
	userEmotesURL   = "https://api.betterttv.net/3/cached/frankerfacez/users/twitch/%s"
	emoteURL        = "https://cdn.frankerfacez.com/emote/%d/4"
)

type Emote struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
}

func (e Emote) ImageURL() string {
	return fmt.Sprintf(emoteURL, e.ID)
}

func GetGlobalEmotes() ([]Emote, error) {
	resp, err := http.Get(globalEmotesURL)
	if err != nil {
		return nil, err
	}

	var emotes []Emote

	err = json.NewDecoder(resp.Body).Decode(&emotes)
	if err != nil {
		return nil, err
	}

	return emotes, nil
}

func GetUserEmotes(userID string) ([]Emote, error) {
	url := fmt.Sprintf(userEmotesURL, userID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var emotes []Emote

	err = json.NewDecoder(resp.Body).Decode(&emotes)
	if err != nil {
		return nil, err
	}

	return emotes, nil
}
