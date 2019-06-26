package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/l-lin/commuting-traffic-info/config"
)

const (
	twitterURL  = "https://api.twitter.com/1.1/statuses/user_timeline.json"
	count       = "5"
	twitterUser = "Ligne%d_RATP"
)

var client = &http.Client{Timeout: time.Second * 10}

// SearchTweets for the given commuting line number
func SearchTweets(result chan *SearchTweetsResult, line int) {
	req, err := http.NewRequest("GET", twitterURL, nil)
	if err != nil {
		result <- &SearchTweetsResult{Tweets: nil, Error: err}
		return
	}
	q := req.URL.Query()
	q.Add("screen_name", fmt.Sprintf(twitterUser, line))
	q.Add("count", count)
	q.Add("tweet_mode", "extended")
	q.Add("exclude_replies", "true")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", "Bearer "+config.GetAPIAuthToken())
	resp, err := client.Do(req)
	if err != nil {
		result <- &SearchTweetsResult{Tweets: nil, Error: err}
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result <- &SearchTweetsResult{Tweets: nil, Error: err}
		return
	}
	if resp.StatusCode != 200 {
		result <- &SearchTweetsResult{Tweets: nil, Error: fmt.Errorf("Could not authenticate. Error status was %d and response body was '%s'", resp.StatusCode, string(body))}
		return
	}
	var data []Tweet
	err = json.Unmarshal(body, &data)
	if err != nil {
		result <- &SearchTweetsResult{Tweets: nil, Error: err}
		return
	}
	result <- &SearchTweetsResult{Tweets: data, Error: nil}
}

// FilterTweets to get only what we need, i.e. no retweet, and today's feed
func FilterTweets(tweets []Tweet) []Tweet {
	filtered := []Tweet{}
	for _, tweet := range tweets {
		creationDate, err := tweet.GetCreationDate()
		if err != nil {
			log.Printf("Could not parse tweet creation date. Error was: %s\n", err.Error())
			continue
		}
		refDate := creationDate.AddDate(0, 0, 1)
		if !tweet.IsRetweet() && time.Now().Before(refDate) {
			filtered = append(filtered, tweet)
		}
	}
	return filtered
}
