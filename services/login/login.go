package login

import (
	"log"

	"hnews/Godeps/_workspace/src/github.com/headzoo/surf"
	"hnews/Godeps/_workspace/src/github.com/headzoo/surf/browser"
)

var (
	loggedInUser map[string]browser.Browser
)

// Login signs the user to Hacker News
func Login(username string, password string) bool {
	bow := surf.NewBrowser()
	err := bow.Open("https://news.ycombinator.com/login")

	// For some reason bow.Form does not find 'form' because bow.Open has not loaded
	// Keep retrying untill
	for {
		if err == nil {
			err = bow.Open("https://news.ycombinator.com/login")
			break
		}
		log.Println(err)
	}

	fm, err := bow.Form("form")

	if err != nil {
		log.Println("Form err:", err)
		return false
	}
	fm.Input("acct", username)
	fm.Input("pw", password)
	err = fm.Submit()
	if err != nil {
		log.Println("Submit form: ", err)
		return false
	}
	if loggedInUser == nil {
		loggedInUser = make(map[string]browser.Browser)
	}
	loggedInUser[username] = *bow
	return true
}

// Upvote upvotes the item with the passed id
func Upvote(id string, username string) bool {
	var bow browser.Browser
	var ok bool
	if bow, ok = loggedInUser[username]; !ok {
		// Not logged in
		return false
	}
	err := bow.Open("https://news.ycombinator.com/item?id=" + id)
	if err != nil {
		// Invalid id
		log.Println(err)
		return false
	}

	err = bow.Click("a#up_" + id)
	if err != nil {
		// Could not click link, already voted?
		log.Println(err)
		return false
	}
	return true
}

// Comment posts a comment on a News item
func Comment(id string, username string, comment string) bool {
	var bow browser.Browser
	var ok bool
	if bow, ok = loggedInUser[username]; !ok {
		// Not logged in
		return false
	}

	err := bow.Open("https://news.ycombinator.com/item?id=" + id)
	if err != nil {
		return false // HN down?
	}

	form := bow.Find("textarea[name='text']")
	s := form.Text()
	log.Println(form, s)
	if err != nil {
		return false
	}

	c := bow.Click("input[type='submit']")
	log.Println(c)

	return true
}

// Reply posts a reply to a comment on a News item
func Reply(id string, username string) bool {

	return true
}
