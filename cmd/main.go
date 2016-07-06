package main

import (
	"flag"
	"hnews/api"
	"hnews/scraper"
	"hnews/services"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	topResource := scraper.Resource{scraper.TopNewsType,
		scraper.TopBaseURL,
		"/top",
		"top", nil}
	askResource := scraper.Resource{scraper.AskNewsType,
		scraper.AskBaseURL,
		"/ask",
		"ask", nil}
	showResource := scraper.Resource{scraper.ShowNewsType,
		scraper.ShowBaseURL,
		"/show",
		"show", nil}
	newestResource := scraper.Resource{scraper.NewestNewsType,
		scraper.NewestBaseURL,
		"/newest",
		"newest", nil}

	// Setup all the scrapers and their ResourceTypes & ResourceURLs
	topScraper := scraper.NewScraper(topResource)
	topResource.BackingStore = topScraper.DatabaseService
	go topScraper.StartScraper(*debug)

	askScraper := scraper.NewScraper(askResource)
	askResource.BackingStore = askScraper.DatabaseService
	go askScraper.StartScraper(*debug)

	showScraper := scraper.NewScraper(showResource)
	showResource.BackingStore = showScraper.DatabaseService
	go showScraper.StartScraper(*debug)

	newestScraper := scraper.NewScraper(newestResource)
	newestResource.BackingStore = newestScraper.DatabaseService
	go newestScraper.StartScraper(*debug)

	// Setup the API by giving it the databases in which the scrapers dumps their data
	api := new(api.API)
	api.Resources = append(api.Resources, topResource)
	api.Resources = append(api.Resources, askResource)
	api.Resources = append(api.Resources, showResource)
	api.Resources = append(api.Resources, newestResource)
	go api.StartAPI(*debug)

	// When closed make sure to call Close on all the underlying bolt.DB instances.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-ch
		topScraper.DatabaseService.Close()
		askScraper.DatabaseService.Close()
		showScraper.DatabaseService.Close()
		newestScraper.DatabaseService.Close()
		services.Commentsdb.Close()
		os.Exit(1)
	}()

	select {} // Block forever and ever
}
