# Generating Tweet Database

Twitter free data for Japanese task of [NTCIR-12 STC](http://ntcir12.noahlab.com.hk/stc.htm) includes only `twitter_id`, so we have to fetch it without money.

I used [GORM](https://github.com/jinzhu/gorm), MySQL and [goquery](https://github.com/PuerkitoBio/goquery).

## preliminary

### Tweet ID Dataset

```
wget https://github.com/mynlp/stc/raw/master/taskdata/twitter_id_str_data.txt.bz2
bzip2 -d twitter_id_str_data.txt.bz2
```

### MySQL

```
mysql> CREATE USER 'trm'@'localhost' IDENTIFIED BY 'trm';
mysql> CREATE DATABASE trm DEFAULT CHARACTER SET utf8;
mysql> GRANT ALL PRIVILEGES ON trm.* TO trm@localhost;
```

If you change username, databasename and password, edit `twitter/database.go`.

### Dependencies

please `go get`

* github.com/PuerkitoBio/goquery
* github.com/go-sql-driver/mysql
* github.com/jinzhu/gorm
* github.comjoho/godotenv

## using

```
go run main.go -reset
```

### Usage

```
go run main.go [-start INT] [-reset] [-check] [-slack STRING]
options(default):
  start(0)     : starting line of dataset
  reset(false) : reset tables of database needed first run
  check(false) : duplication id check of tweet at start line
  slack("")    : channel of slack if notification needed
```

### Log

Please, look at `trm.log`.

### Notification Finish

You can be notified when finished via slack channel.
e.g. `go run main.go -slack general`

Please set `INCOMMING_URL` of environment variable with ["github.comjoho/godotenv"](https://github.com/joho/godotenv).
So, put `.env` file.

.env sample:

```
INCOMMING_URL="https://hooks.slack.com/services/..."
```

## notice

There are packages under subfolders, but not buildable at root.

So, please this:

```
go get -d github.com/sosuke-k/twitter-responding-machine
```
