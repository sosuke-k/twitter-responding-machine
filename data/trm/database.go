package trm

import (
	"fmt"
	"os"
	"time"

	"github.com/codingneo/twittergo"
	// use mysql
	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

// DbName is the database name
const DbName = "trm"

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
func SaveTweet(db *gorm.DB, client *twittergo.Client, id string, update bool) (tweetID int, err error) {
	logger := GetLogger()
	var (
		data  Tweet
		tweet *twittergo.Tweet
	)

	if db.Where("twitter_id = ?", id).First(&data).RecordNotFound() {
		logger.Printf("record(id:%s) not found\n", id)
		update = false
		tweet, err = GetTweet(client, id)
	} else {
		if update && data.Success == 0 {
			logger.Printf("update record(id:%s)\n", id)
			tweet, err = GetTweet(client, id)
		}
		return
	}

	// When there is Twitter API Limit, return.
	if err != nil {
		switch err.(*TwitterError).Op {
		case OpRequest, OpNetwork:
			logger.Fatalln(err)
			os.Exit(1)
		case OpLimit:
			return
		case OpResponse:
			// due to neither API Limit or Network Error
			// e.g. Authorization Error, so on...
			logger.Println("Could not get tweet:" + id)
			logger.Fatalln(err)
			data = Tweet{
				TwitterID: id,
				Success:   0,
			}
		default:
			logger.Fatalln(err)
			panic(err)
		}
	} else {
		name := tweet.User().Name()
		nickname := tweet.User().ScreenName()
		user := User{
			Name:     name,
			Nickname: nickname,
		}
		err = createOrUpdateUser(db, &user)
		if err != nil {
			return
		}
		data = Tweet{
			TwitterID: id,
			Success:   1,
			UserID:    user.ID,
			Text:      tweet.Text(),
			CreatedAt: tweet.CreatedAt(),
		}
	}

	if update {
		err = db.Save(&data).Error
	} else {
		err = db.Create(&data).Error
	}
	if err != nil {
		logger.Printf("could not insert tweet(id:%s)\n", id)
		return
	}
	logger.Printf("insert tweet(id:%s)\n", id)
	tweetID = data.ID
	return
}

func createOrUpdateUser(db *gorm.DB, user *User) (err error) {
	logger := GetLogger()
	if db.Where("nickname = ?", user.Nickname).First(&user).RecordNotFound() {
		err = db.Create(&user).Error
		if err != nil {
			logger.Printf("could not insert user(%s)\n", user.Nickname)
			return
		}
		logger.Printf("insert user(%s)\n", user.Nickname)
	} else {
		logger.Printf("user whose name is %s exists.\n", user.Nickname)
	}
	return
}

// Reset tables of database
func Reset() {
	db, err := DB()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.DropTableIfExists(&User{}).Error
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot drop users table")
		return
	}
	err = db.DropTableIfExists(&Tweet{}).Error
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot drop users table")
		return
	}
	err = db.DropTableIfExists(&Conversation{}).Error
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot drop users table")
		return
	}

	err = db.AutoMigrate(&User{}, &Tweet{}, &Conversation{}).Error
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot migrate tables")
		os.Exit(1)
	}
}
