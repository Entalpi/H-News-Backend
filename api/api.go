package api

// TODO: Start documenting the error codes to the frontend.

import (
	"hnews/services"
	"log"
	"net/http"
	"os"
	"strconv"

	"hnews/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

const (
	debug = true
)

// StartAPI sets up the API and starts it on Heroku port or :8080
func StartAPI() {
	r := gin.Default()
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// GET News from index :from: to index :to:
	r.GET("/v1/news", func(c *gin.Context) {
		from, err0 := strconv.Atoi(c.Query("from"))
		to, err1 := strconv.Atoi(c.Query("to"))
		if err0 != nil || err1 != nil || from <= 0 {
			c.String(http.StatusBadRequest, "Bad index")
			return
		}

		news := services.ReadNews(from, to)
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
