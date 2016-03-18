package main

import (
	"hnews/api"
	"hnews/scraper"
	"hnews/services"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-ch
		services.Close()
		os.Exit(1)
	}()

	go scraper.StartScraper()
	go api.StartAPI()
	select {} // Block forever and ever
}
