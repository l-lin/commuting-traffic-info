package format

import (
	"fmt"

	"github.com/kyokomi/emoji"
	"github.com/l-lin/commuting-traffic-info/traffic"
	"github.com/l-lin/commuting-traffic-info/twitter"
	"github.com/logrusorgru/aurora"
)

// ConsoleFormatter displays output in console
type ConsoleFormatter struct {
}

// Format the output in console
func (f *ConsoleFormatter) Format(lineNb int, s *traffic.Status, tweets []twitter.Tweet) {
	fmt.Printf("%sCommuting traffic for line %d %s\n\n", emoji.Sprint(":train:"), aurora.BrightBlue(lineNb), emoji.Sprint(":train:"))
	fmt.Printf("\t%s\n\n", s)
	if tweets != nil {
		for _, tweet := range tweets {
			fmt.Println(tweet.Render())
		}
	}
}
