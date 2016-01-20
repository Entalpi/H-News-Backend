package api

// TODO: Start documenting the error codes to the frontend.

import (
	"hnews/services"
	"hnews/services/login"
	"log"
	"net/http"
	"os"
	"strconv"

	"hnews/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// StartAPI sets up the API and starts it on Heroku port or :8080
func StartAPI() {
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

	// Given a username and password tries to login with the user and
	// save the session
	r.POST("/v1/login", func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("password")
		success := login.Login(username, password)
		if success {
			c.JSON(http.StatusOK, gin.H{"success": success})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"success": success})
		}
	})

	// Upvotes a specific item, comment, news, etc
	r.POST("/v1/upvote", func(c *gin.Context) {
		newsid := c.Query("id")
		username := c.Query("username")
		success := login.Upvote(newsid, username)
		if success {
			c.JSON(http.StatusOK, gin.H{"success": success})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"success": success})
		}
	})

	// Comment on a news item
	r.POST("/v1/news/comment", func(c *gin.Context) {
		username := c.Query("username")
		newsid := c.Query("id")
		comment := c.Query("comment")
		success := login.Comment(newsid, username, comment)
		if success {
			c.JSON(http.StatusOK, gin.H{"success": success})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"success": success})
		}
	})

	// Post a reply on a comment to a story
	r.POST("/v1/comments/reply", func(c *gin.Context) {
		username := c.Query("username")
		commentid := c.Query("id")
		success := login.Reply(commentid, username)
		if success {
			c.JSON(http.StatusOK, gin.H{"success": success})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"success": success})
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
