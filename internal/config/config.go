package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type App struct {
	Debug    bool
	Channels []*Channel
}

func Read(path string) (*App, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var cfg App

	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

type Channel struct {
	Name string
	Chat *Chat
}

type Chat struct {
	Rewards     []*Trigger
	Commands    []*Trigger
	Spam        Settings
	Middlewares []*Feature
}

type Trigger struct {
	Key         string
	Action      *Feature
	Middlewares []*Feature
}

type Feature struct {
	Type     string
	Settings Settings
}

type Settings map[string]interface{}

func (s Settings) MustString(name string) string {
	v, ok := s[name]
	if !ok {
		log.Fatalf("value not found: %v\n", name)
	}

	str, ok := v.(string)
	if !ok {
		log.Fatalf("value has incorrect type: %t\n", v)
	}

	return str
}

func (s Settings) String(name string) (str string, ok bool) {
	v, ok := s[name]
	if !ok {
		return
	}

	str, ok = v.(string)
	if !ok {
		log.Fatalf("value has incorrect type: %t\n", v)
	}

	return
}

func (s Settings) StringWithDefault(name, defaultVal string) string {
	v, ok := s[name]
	if !ok {
		return defaultVal
	}

	str, ok := v.(string)
	if !ok {
		log.Fatalf("value has incorrect type: %t\n", v)
	}

	return str
}

func (s Settings) MustDuration(name string) time.Duration {
	str := s.MustString(name)

	d, err := time.ParseDuration(str)
	if err != nil {
		log.Fatal(err)
	}

	return d
}

func (s Settings) MustInt(name string) int {
	v, ok := s[name]
	if !ok {
		log.Fatalf("value not found: %v\n", name)
	}

	n, ok := v.(int)
	if !ok {
		log.Fatalf("value has incorrect type: %t\n", v)
	}

	return n
}
