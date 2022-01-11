package bttv

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	globalEmotesURL = "https://api.betterttv.net/3/cached/emotes/global"
	userEmotesURL   = "https://api.betterttv.net/3/cached/users/twitch/%s"
	emoteURL        = "https://cdn.betterttv.net/emote/%s/3x"
)

type Emote struct {
	ID   string `json:"id"`
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

type UserEmotes struct {
	ChannelEmotes []Emote `json:"channelEmotes"`
	SharedEmotes  []Emote `json:"sharedEmotes"`
}

func GetUserEmotes(userID string) (*UserEmotes, error) {
	url := fmt.Sprintf(userEmotesURL, userID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var userEmotes UserEmotes

	err = json.NewDecoder(resp.Body).Decode(&userEmotes)
	if err != nil {
		return nil, err
	}

	return &userEmotes, nil
}
