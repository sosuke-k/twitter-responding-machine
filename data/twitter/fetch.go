/*
Package twitter ...
*/
package twitter

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
)

// Tweet ...
type Tweet struct {
	ID         int
	Success    int    `sql:"not null"`
	ItemID     string `sql:"not null;index"`
	ReplyTo    string `sql:"index"`
	ScreenName string `sql:"index"`
	Name       string
	Time       string
	Text       string
	Replies    []Tweet
}

// Fetch tweet by self id
func (tweet *Tweet) Fetch() (err error) {
	// var b bool
	if tweet.ItemID == "" {
		err = errors.New("This tweet does not have ID attribute.")
		return
	}
	url := "https://twitter.com/statuses/" + tweet.ItemID
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return
	}
	err = tweet.Parse(doc.Find(".permalink-tweet-container .tweet"))
	if err != nil {
		return
	}
	tweet.Replies = []Tweet{}
	doc.Find(".permalink-replies").Find(".stream-item").Each(func(i int, s *goquery.Selection) {
		reply := Tweet{}
		reply.Parse(s.Find(".tweet"))
		reply.ReplyTo = tweet.ItemID
		tweet.Replies = append(tweet.Replies, reply)
	})

	return
}

// Parse from div.tweet
func (tweet *Tweet) Parse(s *goquery.Selection) (err error) {
	var (
		b0 bool
		b1 bool
		b2 bool
		b3 bool
	)
	tweet.ItemID, b0 = s.Attr("data-item-id")
	tweet.ScreenName, b1 = s.Attr("data-screen-name")
	tweet.Name, b2 = s.Attr("data-name")
	tweet.Time, b3 = s.Find("._timestamp").Attr("data-time")
	if !b0 || !b1 || !b2 || !b3 {
		tweet.Success = 0
		err = errors.New("Tweet Parse Error")
		return
	}
	tweet.Success = 1
	tweet.Text = s.Find(".tweet-text").Text()
	return
}
