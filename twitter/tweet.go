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
	return fmt.Sprintf(`┌────────────────────────────────────────────┐
│%s              │
└────────────────────────────────────────────┘
 %s`, t.CreatedAt, t.FullText)
}

// SearchTweetsResult is the result from fetching tweets
type SearchTweetsResult struct {
	Tweets []Tweet
	Error  error
}
