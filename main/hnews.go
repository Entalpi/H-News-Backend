package main

import (
	"hnews/api"
	"hnews/scraper"

	"github.com/davecheney/profile"
)

const (
	// DEBUG Enables extra print outs and information.
	DEBUG = false
)

func main() {
	if DEBUG {
		defer profile.Start(profile.CPUProfile).Stop()
	}
	go scraper.NewScraper()
	go api.NewAPI()
	select {} // Block forever
}
