package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Video struct {
	// json tag to de-serialize json body
	Url string `json:"url" binding:"required"`
}

func main() {
	router := gin.Default()

	router.GET("/private/readiness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	router.GET("/private/liveness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	router.POST("/video", func(context *gin.Context) {
		video := Video{}
		// using BindJson method to serialize body with struct
		if err := context.BindJSON(&video); err != nil {
			context.AbortWithError(http.StatusBadRequest, err)
			return
		}

		fmt.Println(video)
		context.JSON(http.StatusAccepted, &video)
	})

	router.Run(":3000") // listen and serve on 127.0.0.1:3000
}
