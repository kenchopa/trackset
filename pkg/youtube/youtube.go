package youtube

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var lock = &sync.Mutex{}

var service *youtube.Service

type Comment struct {
	Id       string
	VideoId  string
	Content  string
	ParentId *string
	Children []Comment
}

func GetClient() *youtube.Service {
	if service == nil {
		lock.Lock()
		defer lock.Unlock()

		developerKey := os.Getenv("YOUTUBE_API_KEY")
		if developerKey == "" {
			log.Fatalf("Please set a YOUTUBE_API_KEY environment variable.")
		}

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

func GetVideoIdFromYoutubeUrl(youtubeUrl string) string {
	if youtubeUrl == "" {
		log.Fatal("You must provide a full youtube url.")
	}

	videoUrl, err := url.ParseRequestURI(youtubeUrl)
	if err != nil {
		log.Fatal("You must provide a valid url.")
		panic(err)
	}

	// Find the first match for the regular expression.
	var re = regexp.MustCompile(`(?mi)^(?:https?:\/\/)?(?:(?:www\.)?youtube\.com\/(?:(?:v\/)|(?:embed\/|watch(?:\/|\?)){1,2}(?:.*v=)?|.*v=)?|(?:www\.)?youtu\.be\/)([A-Za-z0-9_\-]+)&?.*$`)
	match := re.FindStringSubmatch(videoUrl.String())

	// Extract the video ID from the first capturing group
	if len(match) < 1 {
		log.Fatal("No video ID found.")
	}

	return match[1]
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

func GetVideoInfo(videoID string) *youtube.Video {
	parts := []string{"id", "snippet", "contentDetails", "statistics"}
	call := GetClient().Videos.List(parts).
		Id(videoID)
	response, err := call.Do()
	handleError(err, "")

	if len(response.Items) == 0 {
		log.Fatalf("No video found with ID %v", videoID)
	}

	return response.Items[0]
}

func GetCommentThreads(videoId string, maxResults int64) []Comment {
	parts := []string{"id", "snippet", "replies"}
	call := GetClient().CommentThreads.List(parts).
		VideoId(videoId).
		MaxResults(maxResults)

	response, err := call.Do()
	handleError(err, "")

	// Iterate through each comment and add it to the comments list.
	comments := []Comment{}
	if response.Items == nil || len(response.Items) == 0 {
		return comments
	}

	for _, item := range response.Items {
		parentComment := Comment{item.Id, videoId, item.Snippet.TopLevelComment.Snippet.TextOriginal, nil, nil}

		// make children from replies
		if item.Replies != nil && item.Replies.Comments != nil && len(item.Replies.Comments) != 0 {
			var children []Comment
			for _, replyComment := range item.Replies.Comments {
				children = append(children, Comment{replyComment.Id, videoId, replyComment.Snippet.TextOriginal, &parentComment.Id, nil})
			}
			parentComment.Children = children
		}

		// add top level comment to comments list
		comments = append(comments, parentComment)
	}

	return comments
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
