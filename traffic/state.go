package traffic

import (
	"strings"

	"github.com/kyokomi/emoji"
	"github.com/l-lin/commuting-traffic-info/twitter"
)

var (
	okKeywords  = []string{"rétabli"}
	nokKeywords = []string{"interrompu", "stationne", "perturbé"}

	// OK is when the status is all good
	OK = &Status{emoji.Sprint(":green_heart:")}
	// NOK when the traffic is not good
	NOK = &Status{emoji.Sprint(":anger:")}
	// Warning when the traffic is unknown
	Warning = &Status{emoji.Sprint(":question:")}
)

// Status of the traffic
type Status struct {
	state string
}

func (s *Status) String() string {
	return s.state
}

// GetStatus from the latest tweets
func GetStatus(tweets []twitter.Tweet) *Status {
	if len(tweets) == 0 {
		return OK
	}
	for _, tweet := range tweets {
		for _, keyword := range okKeywords {
			if strings.Contains(tweet.FullText, keyword) {
				return OK
			}
		}
		for _, keyword := range nokKeywords {
			if strings.Contains(tweet.FullText, keyword) {
				return NOK
			}
		}
	}
	return Warning
}
