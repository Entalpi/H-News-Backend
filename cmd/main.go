package main

import (
	"flag"
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

	debug := flag.Bool("debug", false, "Debug mode, defaults to false.")
	flag.Parse()

	go scraper.StartScraper(*debug)
	go api.StartAPI(*debug)
	select {} // Block forever and ever
}
