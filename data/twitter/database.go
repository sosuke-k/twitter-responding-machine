package twitter

import (
	// use mysql
	"fmt"
	"os"

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

// Save tweet to database
func (tweet *Tweet) Save(db *gorm.DB) (err error) {
	err = db.Create(tweet).Error
	if err != nil {
		return
	}
	for i := range tweet.Replies {
		err = db.Create(&tweet.Replies[i]).Error
		if err != nil {
			return
		}
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

	err = db.DropTableIfExists(&Tweet{}).Error
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot drop tweets table")
		os.Exit(1)
	}

	err = db.AutoMigrate(&Tweet{}).Error
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot migrate tables")
		os.Exit(1)
	}
}
