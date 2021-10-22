package creds

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var storage map[string]string

const (
	keyTwitchUser = "twitchuser"
	keyTwitchPass = "twitchpass"
	keyRiotAPIKey = "riotapikey"
)

func Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&storage)
	return err
}

func getValue(key string) string {
	v, ok := storage[key]
	if !ok {
		log.Fatalf("cred value not found: %s\n", key)
	}

	return v
}

func TwitchUser() string {
	return getValue(keyTwitchUser)
}

func TwitchPass() string {
	return getValue(keyTwitchPass)
}

func RiotAPIKey() string {
	return getValue(keyRiotAPIKey)
}
