package main

import (
	"flag"
	"log"

	"github.com/ihrk/microbot/internal/app"
)

func main() {
	var (
		configPath string
		credsPath  string
	)

	flag.StringVar(&configPath, "config", "./config.yml", "")
	flag.StringVar(&credsPath, "creds", "./creds.yml", "")

	flag.Parse()

	log.Fatal(app.LoadConfigAndRun(configPath, credsPath))
}
