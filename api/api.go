package api

import (
	"hnews/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"hnews/services"
	"net/http"
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
		for _, aNews := range news {
			c.JSON(http.StatusOK, aNews)
		}
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
