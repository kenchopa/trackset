package main

import (
	"flag"
	"fmt"
	"log"

	yt "github.com/kenchopa/trackset/pkg/youtube"
)

var (
	query      = flag.String("query", "kenchopa", "Search term")
	maxResults = flag.Int64("max-results", 25, "Max YouTube results")
)

func main() {
	flag.Parse()

	q := *query
	limit := *maxResults
	videos, channels, playlists := yt.Search(q, limit)

	printIDs("Videos", videos)
	printIDs("Channels", channels)
	printIDs("Playlists", playlists)
}

// Print the ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section name "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
