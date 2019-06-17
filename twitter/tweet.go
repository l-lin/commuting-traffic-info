package twitter

import (
	"fmt"
	"time"
)

const searchTweetsURL = "https://api.twitter.com/1.1/search/tweets.json"

// Tweet hold information about tweet
type Tweet struct {
	CreatedAt       string      `json:"created_at"`
	FullText        string      `json:"full_text"`
	RetweetedStatus interface{} `json:"retweeted_status,omitempty"`
}

// IsRetweet returns true if the object "retweeted_statuses" is present in the response
func (t *Tweet) IsRetweet() bool {
	return t.RetweetedStatus != nil
}

// GetCreationDate in time type
func (t *Tweet) GetCreationDate() (time.Time, error) {
	return time.Parse(time.RubyDate, t.CreatedAt)
}

// Render tweet in a pretty format
func (t *Tweet) Render() string {
	size := 45
	fullText := []rune(t.FullText)
	first := t.FullText
	second := ""
	third := ""
	if len(fullText) > size {
		first = string(fullText[0 : size-1])
		if len(fullText) <= size*2 {
			second = string(fullText[size:])
		} else {
			second = string(fullText[size : size*2-1])
			third = string(fullText[size*2-1:])
		}
	}
	return fmt.Sprintf(`
┌────────────────────────────────────────────┐
│%s              │
└────────────────────────────────────────────┘
 %s
 %s
 %s
	`, t.CreatedAt, first, second, third)
}

// SearchTweetsResult is the result from fetching tweets
type SearchTweetsResult struct {
	Tweets []Tweet
	Error  error
}
