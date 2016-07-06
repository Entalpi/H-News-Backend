package api

// TODO: Start documenting the error codes to the frontend.
// TODO: Disable debug mode in Sinatra

import (
	"hnews/scraper"
	"hnews/services"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// API has a pointer to each of the resource DatabaseServices
type API struct {
	TopEndpoint    scraper.Resource
	AskEndpoint    scraper.Resource
	ShowEndpoint   scraper.Resource
	NewestEndpoint scraper.Resource
}

// StartAPI sets up the API and starts it on Heroku port or :8080
func (api *API) StartAPI(debug bool) {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// GET TOP posts from index :from: to index :to:
	r.GET("/v1/top", func(c *gin.Context) {
		from, err0 := strconv.Atoi(c.Query("from"))
		to, err1 := strconv.Atoi(c.Query("to"))
		if err0 != nil || err1 != nil || from <= 0 {
			c.String(http.StatusBadRequest, "Bad index")
			return
		}

		news := api.TopEndpoint.BackingStore.ReadNews(from, to)
		c.JSON(http.StatusOK, gin.H{"values": news})
	})

	// GET ASK posts from index :from: to index :to:
	r.GET("/v1/ask", func(c *gin.Context) {
		from, err0 := strconv.Atoi(c.Query("from"))
		to, err1 := strconv.Atoi(c.Query("to"))
		if err0 != nil || err1 != nil || from <= 0 {
			c.String(http.StatusBadRequest, "Bad index")
			return
		}

		news := api.AskEndpoint.BackingStore.ReadNews(from, to)
		c.JSON(http.StatusOK, gin.H{"values": news})
	})

	// GET SHOw.URL posts from index :from: to index :to:
	r.GET("/v1/show", func(c *gin.Context) {
		from, err0 := strconv.Atoi(c.Query("from"))
		to, err1 := strconv.Atoi(c.Query("to"))
		if err0 != nil || err1 != nil || from <= 0 {
			c.String(http.StatusBadRequest, "Bad index")
			return
		}

		news := api.ShowEndpoint.BackingStore.ReadNews(from, to)
		c.JSON(http.StatusOK, gin.H{"values": news})
	})

	// GET NEWEST posts from index :from: to index :to:
	r.GET("/v1/", func(c *gin.Context) {
		from, err0 := strconv.Atoi(c.Query("from"))
		to, err1 := strconv.Atoi(c.Query("to"))
		if err0 != nil || err1 != nil || from <= 0 {
			c.String(http.StatusBadRequest, "Bad index")
			return
		}

		news := api.NewestEndpoint.BackingStore.ReadNews(from, to)
		c.JSON(http.StatusOK, gin.H{"values": news})
	})

	/** Comment Endpoint **/
	// Gives the comments from a i to j given the provided news id.
	r.GET("/v1/comments", func(c *gin.Context) {
		from, err0 := strconv.Atoi(c.Query("from"))
		to, err1 := strconv.Atoi(c.Query("to"))
		if err0 != nil || err1 != nil || from <= 0 {
			c.String(http.StatusBadRequest, "Bad index")
			return
		}

		id, err := strconv.Atoi(c.Query("newsid"))
		if err != nil {
			c.String(http.StatusBadRequest, "Not a valid item id")
			return
		}

		comments := services.ReadComments(id, from, to)
		c.JSON(http.StatusOK, gin.H{"values": comments})
	})

	/** Login wrapper for login-service **/
	r.POST("/v1/login", func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("password")

		url := "http://localhost:3000/v1/login"
		req, _ := http.NewRequest("POST", url, nil)

		values := req.URL.Query()
		values.Add("username", username)
		values.Add("password", password)
		req.URL.RawQuery = values.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			content, _ := ioutil.ReadAll(resp.Body)
			c.JSON(resp.StatusCode, string(content))
		}
	})

	/** Entry wrapper for login-service **/
	r.POST("/v1/login/entry/upvote", func(c *gin.Context) {
		id := c.Query("id")
		apikey := c.Query("apikey")

		url := "http://localhost:3000/v1/login/entry/upvote"
		req, _ := http.NewRequest("POST", url, nil)

		values := req.URL.Query()
		values.Add("id", id)
		values.Add("apikey", apikey)
		req.URL.RawQuery = values.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			content, _ := ioutil.ReadAll(resp.Body)
			c.JSON(resp.StatusCode, string(content))
		}
	})

	r.POST("/v1/login/entry/comment", func(c *gin.Context) {
		id := c.Query("id")
		comment := c.Query("comment")
		apikey := c.Query("apikey")

		url := "http://localhost:3000/v1/login/entry/comment"
		req, _ := http.NewRequest("POST", url, nil)

		values := req.URL.Query()
		values.Add("id", id)
		values.Add("comment", comment)
		values.Add("apikey", apikey)
		req.URL.RawQuery = values.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			content, _ := ioutil.ReadAll(resp.Body)
			c.JSON(resp.StatusCode, string(content))
		}
	})

	/** Comment wrapper for login-service **/
	r.POST("/v1/login/comment/upvote", func(c *gin.Context) {
		id := c.Query("id")
		apikey := c.Query("apikey")

		url := "http://localhost:3000/v1/login/comment/upvote"
		req, _ := http.NewRequest("POST", url, nil)

		values := req.URL.Query()
		values.Add("id", id)
		values.Add("apikey", apikey)
		req.URL.RawQuery = values.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			content, _ := ioutil.ReadAll(resp.Body)
			c.JSON(resp.StatusCode, string(content))
		}
	})

	r.POST("/v1/login/commment/reply", func(c *gin.Context) {
		id := c.Query("id")
		reply := c.Query("reply")
		apikey := c.Query("apikey")

		url := "http://localhost:3000/v1/login/comment/reply"
		req, _ := http.NewRequest("POST", url, nil)

		values := req.URL.Query()
		values.Add("id", id)
		values.Add("reply", reply)
		values.Add("apikey", apikey)
		req.URL.RawQuery = values.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			content, _ := ioutil.ReadAll(resp.Body)
			c.JSON(resp.StatusCode, string(content))
		}
	})

	r.Run(":" + getPort()) // listen and serve on 0.0.0.0:8080
}

// Tries to get Heroku port otherwise return default 8080
func getPort() string {
	port := os.Getenv("PORT")
	log.Println(port)
	if port != "" {
		return port
	}
	return "8080"
}
