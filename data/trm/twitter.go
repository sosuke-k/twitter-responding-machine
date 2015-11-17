package trm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/codingneo/twittergo"
	"github.com/kurrik/oauth1a"
)

// LoadCredentials load developer information from CREDENTIALS
func LoadCredentials() (client *twittergo.Client, err error) {
	credentials, err := ioutil.ReadFile("CREDENTIALS")
	if err != nil {
		return
	}
	lines := strings.Split(string(credentials), "\n")
	config := &oauth1a.ClientConfig{
		ConsumerKey:    lines[0],
		ConsumerSecret: lines[1],
	}
	user := oauth1a.NewAuthorizedConfig(lines[2], lines[3])
	client = twittergo.NewClient(config, user, "api.twitter.com")
	return
}

// GetTweet return Tweet and error
func GetTweet(client *twittergo.Client, id string) (tweet *twittergo.Tweet, err error) {
	var (
		req  *http.Request
		resp *twittergo.APIResponse
	)

	// parse query
	query := url.Values{}
	query.Set("id", id)
	url := fmt.Sprintf("%v?%v", "/1.1/statuses/show.json", query.Encode())
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		err = &TwitterError{Op: OpRequest, ID: id, Err: err}
		return
	}

	// send request
	resp, err = client.SendRequest(req)
	if err != nil {
		err = &TwitterError{Op: OpNetwork, ID: id, Err: err}
		return
	}

	// parse response
	tweet = &twittergo.Tweet{}
	if e := resp.Parse(tweet); e != nil {
		if rle, ok := e.(twittergo.RateLimitError); ok {
			err = &TwitterError{Op: OpLimit, ID: id, Reset: rle.Reset, Err: e}
			return
		} else if errs, ok := e.(twittergo.Errors); ok {
			var s string
			for i, val := range errs.Errors() {
				s = fmt.Sprintf("Error #%v - ", i+1)
				s += fmt.Sprintf("Code: %v ", val.Code())
				s += fmt.Sprintf("Msg: %v\n", val.Message())
			}
			err = &TwitterError{Op: OpResponse, ID: id, Err: errors.New(s)}
			return
		} else {
			err = &TwitterError{Op: OpResponse, ID: id, Err: e}
			return
		}
	}
	return
}
