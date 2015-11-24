package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sosuke-k/twitter-responding-machine/data/trm/twitter"
)

var reset *bool
var update *bool
var start *int
var startIdx int

// EverySeventeen called per 17 minutes
// func EverySeventeen() {
// 	logger := GetLogger()
//
// 	var (
// 		err    error
// 		db     gorm.DB
// 		client *twittergo.Client
// 	)
//
// 	db, err = DB()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Could not open database: %v\n", err)
// 		os.Exit(1)
// 	}
// 	defer db.Close()
//
// 	client, err = LoadCredentials()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Could not parse CREDENTIALS file: %v\n", err)
// 		os.Exit(1)
// 	}
//
// 	data, err := ioutil.ReadFile("twitter_id_str_data.txt")
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Could not read twitter_id_str_data.txt: %v\n", err)
// 		os.Exit(1)
// 	}
// 	lines := strings.Split(string(data), "\n")
// 	for i := startIdx; i < len(lines); i++ {
// 		logger.Printf("index of lines is %d\n", i)
// 		line := lines[i]
// 		ids := strings.Split(line, "\t")
// 		firstID, err := SaveTweet(&db, client, ids[0], *update)
// 		if err != nil {
// 			if err.(*TwitterError).Op == OpLimit {
// 				startIdx = i
// 				logger.Printf("Next, start at line index %d\n", startIdx)
// 				return
// 			}
// 			panic(err)
// 		}
// 		secondID, err := SaveTweet(&db, client, ids[1], *update)
// 		if err != nil {
// 			if err.(*TwitterError).Op == OpLimit {
// 				startIdx = i
// 				logger.Printf("Next, start at line index %d\n", startIdx)
// 				return
// 			}
// 			panic(err)
// 		}
// 		conversation := Conversation{
// 			FirstTweetID:  firstID,
// 			SecondTweetID: secondID,
// 		}
// 		err = db.Create(&conversation).Error
// 		if err != nil {
// 			panic(err)
// 		}
// 		logger.Println("insert conversation(" + ids[0] + ", " + ids[1] + ")")
// 	}
// }

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
		fmt.Println("db")
		fmt.Println(err.Error())
	}
	defer db.Close()

	var tweet twitter.Tweet
	if db.Where("item_id = ?", itemID).First(&tweet).RecordNotFound() {
		fmt.Println("tweet(id:" + itemID + ") not exists")
		fmt.Println("this start index is ok")
	} else {
		fmt.Println("tweet(id:" + itemID + ") exists")
	}

}

func scrape(start int) {
	db, err := twitter.DB()
	if err != nil {
		fmt.Println("db")
		fmt.Println(err.Error())
	}
	defer db.Close()

	data, err := ioutil.ReadFile("twitter_id_str_data.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read twitter_id_str_data.txt: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	for i := start; i < len(lines); i++ {
		fmt.Printf("index of lines is %d\n", i)
		line := lines[i]
		ids := strings.Split(line, "\t")
		tweet := twitter.Tweet{ItemID: ids[0]}
		err := tweet.Fetch()
		if err != nil {
			fmt.Println("fetch")
			fmt.Println(err.Error())
		}
		err = tweet.Save(&db)
		if err != nil {
			fmt.Println("save")
			fmt.Println(err.Error())
		}
	}
}

func main() {
	// logger := GetLogger()

	check := flag.Bool("check", false, "duplication id check of tweet at start line")
	reset := flag.Bool("reset", false, "reset database")
	// update := flag.Bool("update", false, "update record")
	start := flag.Int("start", 0, "start index")
	flag.Parse()

	if *check {
		duplicationCheck(*start)
		return
	}

	if *reset {
		fmt.Fprintln(os.Stdout, "reset tables...")
		// logger.Println("reset tables...")
		twitter.Reset()
		fmt.Fprintln(os.Stdout, "done")
		// logger.Println("done")
	}

	scrape(*start)
	// startIdx = *start

	// if err := godotenv.Load("../slack/.env"); err != nil {
	// 	fmt.Fprintf(os.Stderr, "%v\n", err)
	// }
	//
	// fmt.Fprintf(os.Stdout, "starting at %d ...\n", *start)
	//
	// c := cron.New()
	// c.AddFunc("@every 17m", EverySeventeen)
	// c.Start()
	//
	// for {
	// 	time.Sleep(10000000000000)
	// 	fmt.Println("still sleeping...")
	// }

}
