package main

import (
	"flag"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	hunt "github.com/kenchopa/trackset/pkg/songhunter"
	yt "github.com/kenchopa/trackset/pkg/youtube"
)

var (
	video = flag.String("video", "", "Video URL to get a tracklist from most used for music set (dj, festival, live concerts...).")
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.Parse()

	videoId := yt.GetVideoIdFromYoutubeUrl(*video)

	comments := yt.GetCommentThreads(videoId, 100)

	songs := []string{}
	for _, comment := range comments {
		songsPerComment := hunt.SearchTrackPattern(comment.Content)
		songs = append(songs, songsPerComment...)

		for _, reply := range comment.Children {
			songsPerReply := hunt.SearchTrackPattern(reply.Content)
			songs = append(songs, songsPerReply...)
		}
	}

	spew.Dump(songs)

	//fmt.Println(re.ReplaceAllString(str, substitution))

	//fmt.Println(re.ReplaceAllString(str, substitution))

	//comments := yt.GetCommentThreads(v, 100)
	//spew.Dump(comments)
}
