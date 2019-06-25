package format

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/l-lin/commuting-traffic-info/traffic"
	"github.com/l-lin/commuting-traffic-info/twitter"
)

// JSONFormatter displays the output in JSON
type JSONFormatter struct {
}

// Format the output in JSON
func (f *JSONFormatter) Format(lineNb int, s *traffic.Status, tweets []twitter.Tweet) {
	output := &jsonOutput{
		Line:   lineNb,
		Status: s.String(),
		Tweets: tweets,
	}
	result, err := json.Marshal(output)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(result))
}

type jsonOutput struct {
	Line   int             `json:"line"`
	Status string          `json:"status"`
	Tweets []twitter.Tweet `json:"tweets"`
}
