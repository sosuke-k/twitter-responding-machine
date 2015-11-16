package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/codingneo/twittergo"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"

	. "./trm"
)

var reset *bool
var update *bool
var start *int
var startIdx int

// EverySeventeen called per 17 minutes
func EverySeventeen() {
	logger := GetLogger()

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

	client, err = LoadCredentials()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse CREDENTIALS file: %v\n", err)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile("twitter_id_str_data.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read twitter_id_str_data.txt: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	for i := startIdx; i < len(lines); i++ {
		logger.Printf("index of lines is %d\n", i)
		line := lines[i]
		ids := strings.Split(line, "\t")
		firstID, err, limit := SaveTweet(&db, client, ids[0], *update)
		if limit {
			startIdx = i
			logger.Printf("Next, start at line index %d\n", startIdx)
			return
		}
		if err != nil {
			panic(err)
		}
		secondID, err, limit := SaveTweet(&db, client, ids[1], *update)
		if limit {
			startIdx = i
			logger.Printf("Next, start at line index %d\n", startIdx)
			return
		}
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
		logger.Println("insert conversation(" + ids[0] + ", " + ids[1] + ")")
	}
}

func main() {
	logger := GetLogger()

	reset = flag.Bool("reset", false, "reset database")
	update = flag.Bool("update", false, "update record")
	start = flag.Int("start", 0, "start index")
	flag.Parse()
	if *reset {
		fmt.Fprintln(os.Stdout, "reset tables...")
		logger.Println("reset tables...")
		Reset()
		fmt.Fprintln(os.Stdout, "done")
		logger.Println("done")
	}
	startIdx = *start

	if err := godotenv.Load("../slack/.env"); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	fmt.Fprintf(os.Stdout, "starting at %d ...\n", *start)

	c := cron.New()
	c.AddFunc("@every 17m", EverySeventeen)
	c.Start()

	for {
		time.Sleep(10000000000000)
		fmt.Println("still sleeping...")
	}

}
