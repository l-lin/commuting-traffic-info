package twitter

const searchTweetsURL = "https://api.twitter.com/1.1/search/tweets.json"

// Tweet hold information about tweet
type Tweet struct {
	CreatedAt       string      `json:"created_at"`
	Text            string      `json:"text"`
	FullText        string      `json:"full_text,omitempty"`
	RetweetedStatus interface{} `json:"retweeted_status,omitempty"`
}

// SearchTweetsResult is the result from fetching tweets
type SearchTweetsResult struct {
	Tweets []Tweet
	Error  error
}
