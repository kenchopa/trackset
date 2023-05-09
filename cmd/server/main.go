package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	hunt "github.com/kenchopa/trackset/pkg/songhunter"
	yt "github.com/kenchopa/trackset/pkg/youtube"
)

type Video struct {
	// json tag to de-serialize json body
	Url string `json:"url" binding:"required"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
			context.JSON(http.StatusOK, gin.H{
				"code":   http.StatusBadRequest,
				"status": "Bad request",
			})
			return
		}

		videoId := yt.GetVideoIdFromYoutubeUrl(video.Url)
		fmt.Println(videoId)

		comments := yt.GetCommentThreads(videoId, 100)
		songs := []string{}
		for _, item := range comments {
			songsPerComment := hunt.SearchTrackPattern(item.Content)
			songs = append(songs, songsPerComment...)
		}

		fmt.Println(video)
		//context.IndentedJSON(http.StatusOK, comments)
		context.JSON(http.StatusOK, gin.H{"data": gin.H{
			"comments": comments,
			"songs":    songs,
		},
			"code":   http.StatusOK,
			"status": "success",
		})
	})

	router.Run(":3000") // listen and serve on 127.0.0.1:3000
}
