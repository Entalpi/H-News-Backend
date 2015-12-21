package scraper

import (
	"hnews/Godeps/_workspace/src/github.com/yhat/scrape"
	"hnews/Godeps/_workspace/src/golang.org/x/net/html"
	"hnews/Godeps/_workspace/src/golang.org/x/net/html/atom"
	"hnews/services"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Scraper struct {
}

func NewScraper() *Scraper {
	scraper := new(Scraper)
	go startScraping()
	return scraper
}

func startScraping() {
	newsCh := make(chan []services.News)
	go scrapeFrontPage(newsCh)

	for {
		select {
		case newNews := <-newsCh:
			go services.SaveNews(newNews)
		}
	}
}

// Parses the front page and sends a []News of the content
func scrapeFrontPage(newsCh chan []services.News) {
	for {
		for i := 1; i <= 16; i++ {
			var news []services.News

			// request and parse the front page
			url := "https://news.ycombinator.com/news?p=" + strconv.Itoa(i)
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
				continue
			}

			root, err := html.Parse(resp.Body)
			if err != nil {
				log.Println(err)
				continue
			}

			points := parsePoints(root)
			ranks := parseRanks(root)
			titles, links := parseArticles(root)
			authors, times, comments := parseSubArticles(root)

			for i := 0; i < len(comments); i++ {
				rank := int32(ranks[i])
				time := times[i]
				title := titles[i]
				link := links[i]

				var author string
				var numPoints int32
				var numComments int32
				if i < len(authors) { // Some stories lack certain properties
					author = authors[i]
					numPoints = int32(points[i])
					numComments = int32(comments[i])
				}

				news = append(news, services.News{rank, title, link, author, numPoints, time, numComments})
			}
			newsCh <- news
		}
		time.Sleep(1 * time.Second)
	}
}

// Parses out the rank of the articles.
func parseRanks(root *html.Node) []int {
	var ranks []int

	rankMatcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Span && n.Parent != nil && n.Parent.Parent != nil {
			grandparent := scrape.Attr(n.Parent.Parent, "class") == "athing"
			self := scrape.Attr(n, "class") == "rank"
			return grandparent && self
		}
		return false
	}

	rankNodes := scrape.FindAll(root, rankMatcher)
	for _, rankNode := range rankNodes {
		text := strings.Replace(scrape.Text(rankNode), ".", "", -1)
		rank, err := strconv.Atoi(text)
		if err != nil {
			rank = 0
		}
		ranks = append(ranks, rank)
	}

	return ranks
}

// Parses out the usernames, timestamps and number of comments.
func parseSubArticles(root *html.Node) ([]string, []time.Time, []int) {
	subarticleMatcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil {
			return scrape.Attr(n.Parent, "class") == "subtext"
		}
		return false
	}

	var authors []string
	var times []time.Time
	var comments []int

	subarticles := scrape.FindAll(root, subarticleMatcher)
	for i := 0; i < len(subarticles); i += 3 {
		authorNode := subarticles[i]
		timeNode := subarticles[i+1]
		commentNode := subarticles[i+2]

		author := scrape.Text(authorNode)

		text := scrape.Text(timeNode)
		time, _ := parseTimeString(text)

		text = scrape.Text(commentNode)
		words := strings.Fields(text)
		numComments, err := strconv.Atoi(words[0])
		if err != nil {
			numComments = 0
		}

		authors = append(authors, author)
		times = append(times, time)
		comments = append(comments, numComments)
	}

	return authors, times, comments
}

// Quantity is of "hours"/"days"/"minutes"
// Text is of "4 hours ago", "41 days ago", etc
func parseTimeString(text string) (time.Time, error) {
	now := time.Now()
	words := strings.Fields(text)

	timeAgo, err := strconv.Atoi(words[0])
	panicOnErr(err)

	var result time.Time
	switch words[1] {
	case "minutes":
		result = now.Add(time.Duration(-timeAgo) * time.Minute)
	case "hours":
		result = now.Add(time.Duration(-timeAgo) * time.Hour)
	case "days":
		result = now.Add(time.Duration(-timeAgo*24) * time.Hour)
	case "day":
		result = now.Add(time.Duration(-timeAgo*24) * time.Hour)
	}
	return result, err
}

func parseArticles(root *html.Node) ([]string, []string) {
	// // define a matcher for main article
	articleMatcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil && n.Parent.Parent != nil {
			return scrape.Attr(n.Parent.Parent, "class") == "athing"
		}
		return false
	}

	var titles []string
	var links []string

	articles := scrape.FindAll(root, articleMatcher)
	for _, article := range articles {
		title := scrape.Text(article)
		link := scrape.Attr(article, "href")
		titles = append(titles, title)
		links = append(links, link)
	}
	return titles, links
}

func parsePoints(root *html.Node) []int {
	// define a matcher for points
	subarticleMatcher := func(n *html.Node) bool {
		// must check for nil values
		if n.DataAtom == atom.Span && n.Parent != nil {
			parent := scrape.Attr(n.Parent, "class") == "subtext"
			return parent
		}
		return false
	}

	var points []int

	pointsNodes := scrape.FindAll(root, subarticleMatcher)
	for _, pointsNode := range pointsNodes {
		pointS := strings.Replace(scrape.Text(pointsNode), " points", "", -1)
		point, err := strconv.Atoi(pointS)
		panicOnErr(err)
		points = append(points, point)
	}
	return points
}

/* Helpers */
func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}