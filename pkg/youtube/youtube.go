package youtube

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var lock = &sync.Mutex{}

const developerKey = "AIzaSyA6uiyrlYQhIyfF4cxkLUSvi47TJ7juWVE"

var service *youtube.Service

func GetClient() *youtube.Service {
	if service == nil {
		lock.Lock()
		defer lock.Unlock()
		fmt.Println("Creating youtube client now.")
		client := &http.Client{
			Transport: &transport.APIKey{Key: developerKey},
		}

		youtubeService, err := youtube.New(client)
		if err != nil {
			log.Fatalf("Error creating new YouTube client: %v", err)
		}
		service = youtubeService
	}

	return service
}

func Search(query string, maxResults int64) (v map[string]string, c map[string]string, p map[string]string) {
	parts := []string{"id", "snippet"}
	call := GetClient().Search.List(parts).
		Q(query).
		MaxResults(maxResults)
	response, err := call.Do()
	handleError(err, "")

	// Group video, channel, and playlist results in separate lists.
	videos := make(map[string]string)
	channels := make(map[string]string)
	playlists := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		case "youtube#channel":
			channels[item.Id.ChannelId] = item.Snippet.Title
		case "youtube#playlist":
			playlists[item.Id.PlaylistId] = item.Snippet.Title
		}
	}

	return videos, channels, playlists
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
