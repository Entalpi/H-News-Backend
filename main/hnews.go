package main

import (
	"hnews/api"
	"hnews/scraper"
)

func main() {
	go scraper.StartScraper()
	go api.StartAPI()
	select {} // Block forever
}
