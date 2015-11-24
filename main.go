package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sosuke-k/twitter-responding-machine/data/trm/twitter"
)

func main() {
	// twitter.Reset()

	// tweet := twitter.Tweet{ItemID: "418033807850496002"}
	// err := tweet.Fetch()
	// if err != nil {
	// 	fmt.Println("fetch")
	// 	fmt.Println(err.Error())
	// }
	db, err := twitter.DB()
	if err != nil {
		fmt.Println("db")
		fmt.Println(err.Error())
	}
	defer db.Close()

	data, err := ioutil.ReadFile("data/twitter_id_str_data.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read twitter_id_str_data.txt: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	for i := 0; i < len(lines); i++ {
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
