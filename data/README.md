# Generating Tweet Database

Twitter free data for Japanese task of [NTCIR-12 STC](http://ntcir12.noahlab.com.hk/stc.htm) includes only `twitter_id`, so we have to fetch it without money.

I used [GORM](https://github.com/jinzhu/gorm), MySQL, [cron](https://github.com/robfig/cron) and [twittergo](https://github.com/kurrik/twittergo) and referred to [twittergo-examples](https://github.com/kurrik/twittergo-examples).

## Structure

![https://gyazo.com/57d73011d2c9135f8a13f7e5759dbef6](https://i.gyazo.com/57d73011d2c9135f8a13f7e5759dbef6.png)

## preliminary

### Tweet ID Dataset

```
wget https://github.com/mynlp/stc/raw/master/taskdata/twitter_id_str_data.txt.bz2
bzip2 -d twitter_id_str_data.txt.bz2
```

### Twitter API

```
touch CREDENTIALS
```

please edit CREDENTIALS

CREDENTIALS sample:

```
<Twitter consumer key>
<Twitter consumer secret>
<Twitter access token>
<Twitter access token secret>
```

### MySQL

```
CREATE USER trm@localhost IDENTIFIED BY trm;
CREATE DATABASE trm DEFAULT CHARACTER SET utf8;
GRANT ALL PRIVILEGES ON trm.* TO trm@localhost;
```

If you change username, databasename and password, edit `trm/database.go`.

## using

```
go run main.go -reset true
```

### Usage

```
go run main.go [-start INT] [-reset BOOL]
options:
  start : starting line of dataset
  reset : reset tables of database
```

### Cron

There is [API Rate Limits](https://dev.twitter.com/rest/public/rate-limiting) of Twitter API Requests.
I call function every fifteen minutes, so many many many long hours......

### Log

Please, look at `trm.log`.

## notice

### Importing local package

Do not put this repository under `GOPATH` folder because of importing local package.
