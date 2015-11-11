package trm

import (
	"fmt"
	"time"

	"github.com/codingneo/twittergo"
	// use mysql
	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

// DbName is the database name
const DbName = "testdatabase"

// UserName is user name of database
const UserName = "trm"

// Password is password of database
const Password = "trm"

// DB returns gorm.DB.
//
// the reference of gorm.DB is there(https://github.com/jinzhu/gorm#query)
func DB() (db gorm.DB, err error) {
	userInfo := UserName + ":" + Password
	dbPath := "@tcp(127.0.0.1:3306)/" + DbName
	dbOption := "?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", userInfo+dbPath+dbOption)
	return
}

// User struct
type User struct {
	ID       int
	Name     string `sql:"not null"`
	Nickname string `sql:"not null;unique_index"`
	Tweets   []Tweet
}

// Tweet struct
type Tweet struct {
	ID        int
	TwitterID string `sql:"not null;unique_index"`
	Success   int    `sql:"not null"`
	UserID    int    `sql:"index"`
	Text      string
	CreatedAt time.Time
}

// Conversation struct
type Conversation struct {
	ID            int
	FirstTweetID  int `sql:"not null;index"`
	SecondTweetID int `sql:"not null;index"`
}

// SaveTweet to database
func SaveTweet(db *gorm.DB, client *twittergo.Client, id string) (tweetID int, err error) {
	tweet, err := GetTweet(client, id)
	var data Tweet
	if err != nil {
		fmt.Println("Could not get tweet:" + id)
		data = Tweet{
			TwitterID: id,
			Success:   0,
		}
	} else {
		//check users table
		var user User
		name := tweet.User().Name()
		nickname := tweet.User().ScreenName()
		if db.Where("nickname = ?", nickname).First(&user).RecordNotFound() {
			user = User{
				Name:     name,
				Nickname: nickname,
			}
			err = db.Create(&user).Error
			if err != nil {
				return
			}
		}

		data = Tweet{
			TwitterID: id,
			Success:   1,
			UserID:    user.ID,
			Text:      tweet.Text(),
			CreatedAt: tweet.CreatedAt(),
		}
	}
	err = db.Create(&data).Error
	if err != nil {
		return
	}
	tweetID = data.ID
	return
}
