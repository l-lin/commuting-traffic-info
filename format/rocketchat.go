package format

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kyokomi/emoji"
	"github.com/l-lin/commuting-traffic-info/traffic"
	"github.com/l-lin/commuting-traffic-info/twitter"
	"github.com/logrusorgru/aurora"
)

var colors = map[int]string{
	1:  "#ffbe00",
	2:  "#0354c3",
	3:  "#6e6e00",
	4:  "#a0006e",
	5:  "#ff5900",
	6:  "#83db73",
	7:  "#fd84b4",
	8:  "#d282bf",
	9:  "#d2d200",
	10: "#db9500",
	11: "#5b230a",
	12: "#04623c",
	13: "#82c6e6",
}

// RocketChatFormatter sends a webhook to rocket.chat
type RocketChatFormatter struct {
	Webhook string
}

// Format the commuting status info by sending a webhook to rocket.chat
func (f *RocketChatFormatter) Format(lineNb int, s *traffic.Status, tweets []twitter.Tweet) {
	msg := &RocketChatMessage{
		IconEmoji: ":train:",
		Text:      fmt.Sprintf("Commuting traffic info for line %d: %s", lineNb, s.String()),
	}
	if tweets != nil && len(tweets) > 0 {
		var buffer bytes.Buffer
		for _, t := range tweets {
			buffer.WriteString(fmt.Sprintf("%s %s\n%s\n\n", emoji.Sprint(":bird:"), t.CreatedAt, t.FullText))
		}
		u := tweets[0].User
		attachment := &Attachment{
			Title:     u.Name,
			TitleLink: u.URL,
			Text:      buffer.String(),
			Color:     colors[lineNb],
		}
		msg.Attachments = []Attachment{*attachment}
	}
	result, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln(aurora.BrightRed(err))
	}
	resp, err := http.Post(f.Webhook, "application/json", strings.NewReader(string(result)))
	if err != nil {
		log.Fatalln(aurora.BrightRed(err))
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalln(aurora.BrightRed("Could not call Rocket.Chat webhook. Status was: " + resp.Status))
	}
}

// RocketChatMessage is the message to send to rocket chat
type RocketChatMessage struct {
	IconEmoji   string       `json:"icon_emoji"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

// Attachment of the rocket.chat message
type Attachment struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
	ImageURL  string `json:"image_url"`
	Color     string `json:"color"`
}
