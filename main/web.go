package main

import (
	"hnews/api"
	"hnews/scraper"
)

func main() {
	go scraper.NewScraper()
	go api.NewAPI()
	select {} // Block forever
}
