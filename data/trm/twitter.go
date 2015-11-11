package trm

import (
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
	query := url.Values{}
	query.Set("id", id)
	url := fmt.Sprintf("%v?%v", "/1.1/statuses/show.json", query.Encode())
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		fmt.Println("Please see https://twitter.com/statuses/" + id)
		return
	}
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		fmt.Println("Please see https://twitter.com/statuses/" + id)
		return
	}
	tweet = &twittergo.Tweet{}
	err = resp.Parse(tweet)
	if err != nil {
		if rle, ok := err.(twittergo.RateLimitError); ok {
			fmt.Printf("Rate limited, reset at %v\n", rle.Reset)
		} else if errs, ok := err.(twittergo.Errors); ok {
			for i, val := range errs.Errors() {
				fmt.Printf("Error #%v - ", i+1)
				fmt.Printf("Code: %v ", val.Code())
				fmt.Printf("Msg: %v\n", val.Message())
			}
		} else {
			fmt.Printf("Problem parsing response: %v\n", err)
		}
		fmt.Println("Please see https://twitter.com/statuses/" + id)
		return
	}
	return
}
