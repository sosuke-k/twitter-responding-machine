package trm

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/codingneo/twittergo"
	"github.com/kurrik/oauth1a"

	"../../slack"
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
func GetTweet(client *twittergo.Client, id string) (tweet *twittergo.Tweet, err error, rateLimit bool) {
	logger := GetLogger()
	rateLimit = false
	var (
		req  *http.Request
		resp *twittergo.APIResponse
	)
	query := url.Values{}
	query.Set("id", id)
	url := fmt.Sprintf("%v?%v", "/1.1/statuses/show.json", query.Encode())
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Printf("Could not parse request: %v\n", err)
		logger.Println("Please see https://twitter.com/statuses/" + id)
		return
	}
	resp, err = client.SendRequest(req)
	if err != nil {
		logger.Printf("Could not send request: %v\n", err)
		logger.Println("Please see https://twitter.com/statuses/" + id)
		err := slack.Post("Could not send request. Please see log file.")
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}
	tweet = &twittergo.Tweet{}
	err = resp.Parse(tweet)
	if err != nil {
		if rle, ok := err.(twittergo.RateLimitError); ok {
			logger.Printf("Rate limited, reset at %v\n", rle.Reset)
			rateLimit = true
			return
		} else if errs, ok := err.(twittergo.Errors); ok {
			for i, val := range errs.Errors() {
				logger.Printf("Error #%v - ", i+1)
				logger.Printf("Code: %v ", val.Code())
				logger.Printf("Msg: %v\n", val.Message())
			}
		} else {
			logger.Printf("Problem parsing response: %v\n", err)
		}
		logger.Println("Please see https://twitter.com/statuses/" + id)
		return
	}
	return
}
