/*
Package twitter ...
*/
package twitter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Operation ...
type Operation struct {
	Query         int
	Request       int
	Authorization int
	NotExisting   int
	Parse         int
}

// Op is alternative Enum object for Operation
var Op = Operation{1, 2, 3, 4, 5}

// A Error records a failed get of tweet.
type Error struct {
	Op  int    // the failing Operation (Query, Request, Authorization, Parse)
	ID  string // the twitter id
	URL string // the definitive url
	Err error  // the reason the get failed
}

func (e *Error) Error() string {
	return " https://twitter.com/statuses/" + e.ID + ": " + e.Err.Error()
}

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
	if tweet.ItemID == "" {
		err = errors.New("Tweet.ItemID is empty")
		err = &Error{Op: Op.Query, Err: err}
		return
	}

	url := "https://twitter.com/statuses/" + tweet.ItemID
	doc, err := goquery.NewDocument(url)
	if err != nil {
		err = &Error{Op: Op.Request, ID: tweet.ItemID, Err: err}
		return
	}

	if !strings.Contains(doc.Url.Path, tweet.ItemID) {
		err = errors.New("May be redirected because of authorization error")
		err = &Error{Op: Op.Authorization, ID: tweet.ItemID, URL: doc.Url.String(), Err: err}
		return
	}

	err = tweet.Parse(doc.Find(".permalink-tweet-container .tweet"))
	if err != nil {
		if notExists := checkExisting(doc); notExists {
			err = errors.New("this page not exists")
			err = &Error{Op: Op.NotExisting, ID: tweet.ItemID, URL: doc.Url.String(), Err: err}
		} else {
			err = &Error{Op: Op.Parse, ID: tweet.ItemID, URL: doc.Url.String(), Err: err}
		}
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
	success := false
	attrs := []string{
		"data-item-id",
		"data-screen-name",
		"data-name",
	}
	data := map[string]string{}

	for _, attr := range attrs {
		var value string
		if value, success = s.Attr(attr); !success {
			tweet.Success = 0
			err = fmt.Errorf("not having %s attribute", attr)
			return
		}
		data[attr] = value
	}

	tweet.ItemID = data["data-item-id"]
	tweet.ScreenName = data["data-screen-name"]
	tweet.Name = data["data-name"]
	tweet.Success = 1

	// if could get the above attribues, allow the following values to be blank.
	tweet.Time, _ = s.Find("._timestamp").Attr("data-time")
	tweet.Text = s.Find(".tweet-text").Text()
	return
}

func checkExisting(doc *goquery.Document) (notExists bool) {
	if doc.Find(".body-content h1").Text() == "Sorry, that page doesn’t exist!" {
		notExists = true
	}
	return
}
