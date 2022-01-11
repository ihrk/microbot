package config

import (
	"log"
	"os"
	"strings"
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
	str, ok := s.String(name)
	if !ok {
		log.Fatalf("value not found: %v\n", name)
	}

	return str
}

func (s Settings) String(name string) (string, bool) {
	v, ok := s[name]
	if !ok {
		return "", false
	}

	str, ok := v.(string)
	if !ok {
		log.Fatalf("value has incorrect type: %t\n", v)
	}

	return str, true
}

func (s Settings) StringFromSet(name string, set []string) (string, bool) {
	str, ok := s.String(name)
	if !ok {
		return "", false
	}

	for _, e := range set {
		if str == e {
			return str, true
		}
	}

	log.Fatalf("unexpected value: %v, expected values: %v\n",
		str, strings.Join(set, ", "))

	return "", false
}

func (s Settings) StringFromSetWithDefault(name string, set []string, d string) string {
	str, ok := s.StringFromSet(name, set)
	if !ok {
		return d
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

func (s Settings) Int(name string) (int, bool) {
	v, ok := s[name]
	if !ok {
		return 0, false
	}

	n, ok := v.(int)
	if !ok {
		log.Fatalf("value has incorrect type: %t\n", v)
	}

	return n, true
}

func (s Settings) MustInt(name string) int {
	n, ok := s.Int(name)
	if !ok {
		log.Fatalf("value not found: %v\n", name)
	}

	return n
}

// Bool returns false if value is not found
func (s Settings) Bool(name string) bool {
	v, ok := s[name]
	if !ok {
		return false
	}

	b, ok := v.(bool)
	if !ok {
		log.Fatalf("value has incorrect type: %t\n", v)
	}

	return b
}

func (s Settings) Strings(name string) []string {
	v, ok := s[name]
	if !ok {
		return nil
	}

	arr, ok := v.([]interface{})
	if !ok {
		log.Fatalf("value has incorrect type: %t", v)
	}

	a := make([]string, len(arr))

	for i := range arr {
		str, ok := arr[i].(string)
		if !ok {
			log.Fatalf("value has incorrect type: %t", arr[i])
		}

		a[i] = str
	}

	return a
}
