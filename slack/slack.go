/*
Package slack has only `Post(message)` function.
This function can post just a message.

Please set `INCOMMING_URL` of environment variable.
*/
package slack

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

type payload struct {
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	Text      string `json:"text"`
}

// Post send a plain message
func Post(text string) error {
	webhookURL := os.Getenv("INCOMMING_URL")

	p, err := json.Marshal(&payload{
		Username:  "trm",
		IconEmoji: ":trm:",
		Text:      text,
	})
	if err != nil {
		return err
	}

	_, err = http.PostForm(webhookURL, url.Values{
		"payload": []string{string(p)},
	})
	if err != nil {
		return err
	}
	return nil
}
