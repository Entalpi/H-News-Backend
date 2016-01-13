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
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		services.Close()
		os.Exit(1)
	}()

	go scraper.StartScraper()
	go api.StartAPI()
	select {} // Block forever
}
