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
	commentsCh := make(chan []services.Comment)
	go scrapeFrontPage(newsCh)
	go scrapeComments(commentsCh)

	for {
		select {
		case newNews := <-newsCh:
			go services.SaveNews(newNews)
		case newComments := <-commentsCh:
			go services.SaveComments(newComments)
		}
	}
}

/********************** News **********************/
// Parses the front page and sends a []News of the content
func scrapeFrontPage(newsCh chan []services.News) {
	for {
		for i := 1; i <= 16; i++ {
			var news []services.News

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
			authors := parseAuthors(root)
			times := parseTimes(root)
			comments := parseNumComments(root)
			ids := parseIDs(root)

			for i := 0; i < len(ranks); i++ {
				rank := int32(ranks[i])
				title := titles[i]

				var time time.Time
				if i < len(times) {
					time = times[i]
				}

				var link string
				if i < len(links) {
					link = links[i]
				}

				var author string
				if i < len(authors) {
					author = authors[i]
				}

				var numPoints int32
				if i < len(points) {
					numPoints = int32(points[i])
				}

				var numComments int32
				if i < len(comments) {
					numComments = int32(comments[i])
				}

				var id int32
				if i < len(ids) {
					id = int32(ids[i])
				}

				news = append(news, services.News{id, rank, title, link, author, numPoints, time, numComments})
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

// Parses the authors of each Story on the frontpage
func parseAuthors(root *html.Node) []string {
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil {
			parent := scrape.Attr(n.Parent, "class") == "subtext"
			href := scrape.Attr(n, "href")
			if len(href) > 4 && href[0:4] == "user" {
				return parent && true
			}
		}
		return false
	}

	var authors []string
	authorNodes := scrape.FindAll(root, matcher)
	for _, authorNode := range authorNodes {
		authors = append(authors, scrape.Text(authorNode))
	}
	return authors
}

// Parses the number of comments on each Story on the frontpage
func parseNumComments(root *html.Node) []int {
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil {
			parent := scrape.Attr(n.Parent, "class") == "subtext"
			href := scrape.Attr(n, "href")
			if len(href) > 4 && href[0:4] == "item" {
				return parent && true
			}
		}
		return false
	}

	var numComms []int // Number of comments for each Story
	commentNodes := scrape.FindAll(root, matcher)
	for _, commentNode := range commentNodes {
		text := scrape.Text(commentNode)
		words := strings.Fields(text)
		numComm, err := strconv.Atoi(words[0])
		if err != nil {
			numComm = 0
		}
		numComms = append(numComms, numComm)
	}
	return numComms
}

// Parses the timestamps from the HTML of the frontpage
func parseTimes(root *html.Node) []time.Time {
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil {
			return scrape.Attr(n.Parent, "class") == "age"
		}
		return false
	}

	var dates []time.Time
	timeNodes := scrape.FindAll(root, matcher)
	for _, timeNode := range timeNodes {
		text := scrape.Text(timeNode)
		date, err := parseTimeString(text)
		if err != nil {
			date = time.Now()
		}
		dates = append(dates, date)
	}
	return dates
}

// Parses the id of each Story on the page
func parseIDs(root *html.Node) []int {
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Span && n.Parent != nil {
			parent := scrape.Attr(n.Parent, "class") == "subtext"
			self := scrape.Attr(n, "class") == "score"
			return parent && self
		}
		return false
	}

	var ids []int
	idNodes := scrape.FindAll(root, matcher)
	for _, idNode := range idNodes {
		var id int
		idTemp := scrape.Attr(idNode, "id")
		if len(idTemp) <= 6 {
			id = 0
			ids = append(ids, id)
			continue
		}
		id, err := strconv.Atoi(idTemp[6:len(idTemp)])
		if err != nil {
			id = 0
		}
		ids = append(ids, id)
	}
	return ids
}

// Quantity is of "hours"/"days"/"minutes"
// Text is of "4 hours ago", "41 days ago", etc
func parseTimeString(text string) (time.Time, error) {
	now := time.Now()
	words := strings.Fields(text)

	timeAgo, err := strconv.Atoi(words[0])
	if err != nil {
		return now, err
	}

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

// Parses the title and the link of each Story on the page
func parseArticles(root *html.Node) ([]string, []string) {
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil && n.Parent.Parent != nil {
			return scrape.Attr(n.Parent.Parent, "class") == "athing"
		}
		return false
	}

	var titles []string
	var links []string

	articles := scrape.FindAll(root, matcher)
	for _, article := range articles {
		title := scrape.Text(article)
		link := scrape.Attr(article, "href")
		titles = append(titles, title)
		links = append(links, link)
	}
	return titles, links
}

func parsePoints(root *html.Node) []int {
	matcher := func(n *html.Node) bool {
		// must check for nil values
		if n.DataAtom == atom.Span && n.Parent != nil {
			parent := scrape.Attr(n.Parent, "class") == "subtext"
			self := scrape.Attr(n, "class") == "score"
			return parent && self
		}
		return false
	}

	var points []int

	pointsNodes := scrape.FindAll(root, matcher)
	for _, pointsNode := range pointsNodes {
		pointS := strings.Replace(scrape.Text(pointsNode), " points", "", -1)
		point, err := strconv.Atoi(pointS)
		if err != nil {
			point = 0
		}
		points = append(points, point)
	}
	return points
}

/********************** News **********************/

/******************** Comments ********************/
// Scrapes the Comments for every News item currently in the database.
func scrapeComments(commentsCh chan []services.Comment) {
	for {
		ids := services.ReadNewsIds()
		for _, id := range ids {
			url := "https://news.ycombinator.com/item?id=" + strconv.Itoa(int(id))

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

			comments := parseComments(root, id)
			commentsCh <- comments
		}
		time.Sleep(1 * time.Second)
	}
}

// Parses all the Comments for a particular News item.
func parseComments(root *html.Node, newsid int32) []services.Comment {
	offsets := parseOffsets(root)
	ids := parseCommentIDs(root)
	authors := parseCommentAuthors(root)
	times := parseCommentTimes(root)
	texts := parseCommentText(root)

	var rootComments []services.Comment
	for i, _ := range authors {
		comment := services.Comment{int32(i), newsid, int32(ids[i]), int32(offsets[i]),
			times[i], authors[i], texts[i]}
		rootComments = append(rootComments, comment)
	}
	return rootComments
}

// Parses the level for each Comment in the comment tree. Interval: 0-inf.
func parseOffsets(root *html.Node) []int {
	offsetMatcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Img && n.Parent != nil {
			parent := scrape.Attr(n.Parent, "class") == "ind"
			return parent
		}
		return false
	}

	offsets := scrape.FindAll(root, offsetMatcher)

	var norm []int
	for _, offset := range offsets {
		lvl, _ := strconv.Atoi(scrape.Attr(offset, "width"))
		norm = append(norm, int(lvl/40))
	}
	return norm
}

// Parses all the comments authors
func parseCommentAuthors(root *html.Node) []string {
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil {
			return scrape.Attr(n.Parent, "class") == "comhead"
		}
		return false
	}

	var authors []string
	authorNodes := scrape.FindAll(root, matcher)
	for _, authorNode := range authorNodes {
		authors = append(authors, scrape.Text(authorNode))
	}
	return authors
}

// Parses all the comments itemids
func parseCommentIDs(root *html.Node) []int {
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil {
			return scrape.Attr(n.Parent, "class") == "age"
		}
		return false
	}

	var ids []int
	idNodes := scrape.FindAll(root, matcher)
	for _, idNode := range idNodes {
		var id int
		href := scrape.Attr(idNode, "href")
		if len(href) <= 8 {
			id = 0
			ids = append(ids, id)
			continue
		}
		id, err := strconv.Atoi(href[8:len(href)])
		if err != nil {
			id = 0
		}
		ids = append(ids, id)
	}
	return ids
}

// Parses all the comments timestamps
func parseCommentTimes(root *html.Node) []time.Time {
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil {
			return scrape.Attr(n.Parent, "class") == "age"
		}
		return false
	}

	var dates []time.Time
	timeNodes := scrape.FindAll(root, matcher)
	for _, timeNode := range timeNodes {
		date, err := parseTimeString(scrape.Text(timeNode))
		if err != nil {
			date = time.Now()
		}
		dates = append(dates, date)
	}
	return dates
}

// Parses the text of all Comments for a News
func parseCommentText(root *html.Node) []string {
	textMatcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Span && n.Parent != nil {
			parent := scrape.Attr(n.Parent, "class") == "comment"
			return parent
		}
		return false
	}

	textNodes := scrape.FindAll(root, textMatcher)

	var texts []string
	for _, text := range textNodes {
		content := scrape.Text(text)
		// BUG: Remove trailing trash from the 'Reply' HTML node ...
		texts = append(texts, content) //[0:len(content)-5])
	}
	return texts
}

/******************** Comments ********************/
