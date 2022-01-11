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

	flag.StringVar(&configPath, "config", "./config.yml", "path to configuration file")
	flag.StringVar(&credsPath, "creds", "./creds.yml", "path to file with creds")

	flag.Parse()

	log.Fatalf("dial failed with error: %v\n",
		app.LoadConfigAndRun(configPath, credsPath))
}
