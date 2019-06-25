package format

import (
	"github.com/l-lin/commuting-traffic-info/traffic"
	"github.com/l-lin/commuting-traffic-info/twitter"
)

// Formatter displays output
type Formatter interface {
	Format(lineNb int, s *traffic.Status, tweets []twitter.Tweet)
}
