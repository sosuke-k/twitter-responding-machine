package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sosuke-k/twitter-responding-machine/data/logger"
	"github.com/sosuke-k/twitter-responding-machine/data/slack"
	"github.com/sosuke-k/twitter-responding-machine/data/twitter"
)

func duplicationCheck(start int) {

	data, err := ioutil.ReadFile("twitter_id_str_data.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read twitter_id_str_data.txt: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	line := lines[start]
	itemID := strings.Split(line, "\t")[0]

	db, err := twitter.DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	var tweet twitter.Tweet
	if db.Where("item_id = ?", itemID).First(&tweet).RecordNotFound() {
		fmt.Fprintf(os.Stdout, "tweet(id:%s) not exists\n", itemID)
		fmt.Println("this start index is ok")
	} else {
		fmt.Fprintf(os.Stdout, "tweet(id:%s) exists\n", itemID)
		fmt.Println("this start index is not good")
	}

}

func gather(start int) {
	logger := logger.GetInstance()

	db, err := twitter.DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	data, err := ioutil.ReadFile("twitter_id_str_data.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open database: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	for i := start; i < len(lines); i++ {
		fmt.Printf("index of lines is %d\n", i)
		logger.Printf("index of lines is %d\n", i)
		line := lines[i]
		ids := strings.Split(line, "\t")
		tweet := twitter.Tweet{ItemID: ids[0]}
		err := tweet.Fetch()
		if err != nil {
			logger.Println("Could not fetch tweet:")
			logger.Printf("    item is = %s\n", tweet.ItemID)
			logger.Println(err.Error())
			fmt.Fprintf(os.Stderr, "Could not fetch tweet: %v\n", err)
		}
		err = tweet.Save(&db)
		if err != nil {
			logger.Println("Could not save tweet:")
			logger.Printf("    item is = %s\n", tweet.ItemID)
			logger.Println(err.Error())
			fmt.Fprintf(os.Stderr, "Could not save tweet: %v\n", err)
		}
		logger.Println("Successed inserting tweet:")
		logger.Printf("    item is = %s\n", tweet.ItemID)
	}
}

func retryGather(start int, channel string) {
	logger := logger.GetInstance()

	db, err := twitter.DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	data, err := ioutil.ReadFile("failed_ids.tsv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open failed_post_ids.tsv: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	fmt.Fprintf(os.Stdout, "%d lines in failed_post_ids.tsv", len(lines))
	for i := start; i < len(lines); i++ {
		fmt.Printf("index of lines is %d\n", i)
		logger.Printf("index of lines is %d\n", i)
		itemID := lines[i]
		tweet := twitter.Tweet{ItemID: itemID}
		err := tweet.Fetch()
		if err != nil {
			logger.Println("Could not fetch tweet:")
			logger.Printf("    item is = %s\n", tweet.ItemID)
			logger.Println(err.Error())
			fmt.Fprintf(os.Stderr, "Could not fetch tweet: %v\n", err)

			op := err.(*twitter.Error).Op
			if op != twitter.Op.Authorization && op != twitter.Op.NotExisting {
				slack.Post(channel, err.Error())
				continue
			}
		}
		err = tweet.Save(&db)
		if err != nil {
			slack.Post(channel, "Could not save tweet:\n  https://twitter.com/statuses/"+tweet.ItemID)
			logger.Println("Could not save tweet:")
			logger.Printf("    item is = %s\n", tweet.ItemID)
			logger.Println(err.Error())
			fmt.Fprintf(os.Stderr, "Could not save tweet: %v\n", err)
		}
		logger.Println("Successed inserting tweet:")
		logger.Printf("    item is = %s\n", tweet.ItemID)
	}

}

func main() {
	logger := logger.GetInstance()

	check := flag.Bool("check", false, "duplication id check of tweet at start line")
	reset := flag.Bool("reset", false, "reset database")
	start := flag.Int("start", 0, "start index")
	retry := flag.Bool("retry", false, "retry failed ids")
	channel := flag.String("slack", "", "channel of slack if notification needed")
	flag.Parse()

	fmt.Fprintf(os.Stdout, "starting line number = %d\n", *start)

	if *check {
		duplicationCheck(*start)
		return
	}

	if *channel != "" {
		if err := godotenv.Load(".env"); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}

	if *retry {
		if *channel == "" {
			fmt.Fprintln(os.Stdout, "Please, give channel argument")
			os.Exit(1)
		}
		retryGather(*start, *channel)
		slack.Post(*channel, "finished retry gathering!")
		return
	}

	if *reset {
		fmt.Fprintln(os.Stdout, "reset tables...")
		logger.Println("reset tables...")
		twitter.Reset()
		fmt.Fprintln(os.Stdout, "done")
		logger.Println("done")
	}

	gather(*start)

	if *channel != "" {
		slack.Post(*channel, "finished gathering!")
	}

}
