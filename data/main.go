package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/codingneo/twittergo"
	"github.com/jinzhu/gorm"

	. "./trm"
)

// ParseTweet parse twittergo.Tweet to trm.Tweet.
func ParseTweet(id string, tweet *twittergo.Tweet) (model Tweet, err error) {
	if tweet == nil {
		model = Tweet{
			TwitterID: id,
			Success:   0,
		}
		return
	}
	// t, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	// if err != nil {
	// 	return
	// }
	text := tweet.Text()
	t := tweet.CreatedAt()
	model = Tweet{
		TwitterID: id,
		Success:   1,
		Text:      text,
		CreatedAt: t,
	}
	return
}

func main() {

	var (
		err    error
		db     gorm.DB
		client *twittergo.Client
	)

	db, err = DB()
	if err != nil {
		return
	}
	defer db.Close()
	db.DropTableIfExists(&User{})
	db.DropTableIfExists(&Tweet{})
	db.DropTableIfExists(&Conversation{})

	err = db.AutoMigrate(&User{}, &Tweet{}, &Conversation{}).Error
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot migrate tables")
		os.Exit(1)
	}

	client, err = LoadCredentials()
	if err != nil {
		fmt.Printf("Could not parse CREDENTIALS file: %v\n", err)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile("twitter_id_str_data.txt")
	if err != nil {
		fmt.Printf("Could not read twitter_id_str_data.txt: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if i < 11 {
			fmt.Println("i = ", strconv.Itoa(i))
			ids := strings.Split(line, "\t")
			firstID, err := SaveTweet(&db, client, ids[0])
			if err != nil {
				panic(err)
			}
			secondID, err := SaveTweet(&db, client, ids[1])
			if err != nil {
				panic(err)
			}
			conversation := Conversation{
				FirstTweetID:  firstID,
				SecondTweetID: secondID,
			}
			err = db.Create(&conversation).Error
			if err != nil {
				panic(err)
			}
		}
	}
}
