package api

import (
	"hnews/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"hnews/services"
	"log"
	"net/http"
	"os"
	"strconv"
)

type API struct {
}

func NewAPI() *API {
	api := new(API)
	go api.StartAPI()
	return api
}

func (*API) StartAPI() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	// GET News from index :from: to index :to:
	r.GET("/v1/news", func(c *gin.Context) {
		from, err0 := strconv.Atoi(c.Query("from"))
		to, err1 := strconv.Atoi(c.Query("to"))
		if err0 != nil || err1 != nil || from <= 0 {
			c.String(http.StatusBadRequest, "Bad index")
			return
		}

		news := services.ReadNews(from, to)
		c.JSON(http.StatusOK, gin.H{"news": news})
	})

	r.GET("/v1/comments", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("newsid"))
		if err != nil {
			c.String(http.StatusBadRequest, "Not a valid item id")
			return
		}

		comments := services.ReadComments(id)
		c.JSON(http.StatusOK, gin.H{"comments": comments})
	})

	r.Run(":" + getPort()) // listen and serve on 0.0.0.0:8080
}

// Tries to get Heroku port otherwise return default 8080
func getPort() string {
	port := os.Getenv("PORT")
	log.Println(port)
	if port != "" {
		return port
	} else {
		return "8080"
	}
}
