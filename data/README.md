# Generating Tweet Database

Twitter free data for Japanese task of [NTCIR-12 STC](http://ntcir12.noahlab.com.hk/stc.htm) includes only `twitter_id`, so we have to fetch it without money.

I used [GORM](https://github.com/jinzhu/gorm), MySQL, [cron](https://github.com/robfig/cron) and [twittergo](https://github.com/kurrik/twittergo) and referred to [twittergo-examples](https://github.com/kurrik/twittergo-examples).

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
mysql> CREATE USER 'trm'@'localhost' IDENTIFIED BY 'trm';
mysql> CREATE DATABASE trm DEFAULT CHARACTER SET utf8;
mysql> GRANT ALL PRIVILEGES ON trm.* TO trm@localhost;
```

If you change username, databasename and password, edit `trm/database.go`.

### Dependencies

please `go get`

* github.com/codingneo/twittergo
* github.com/go-sql-driver/mysql
* github.com/jinzhu/gorm
* github.com/kurrik/oauth1a
* github.com/robfig/cron

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
I call function every 17 minutes, so many many many long hours......

First, I try every 15 minutes, but this

```
2015/11/12 13:01:41 Rate limited, reset at 2015-11-12 13:15:00 +0900 JST
2015/11/12 13:01:41 Next, start at line index 90
2015/11/12 13:15:00 index of lines is 90
2015/11/12 13:15:00 Rate limited, reset at 2015-11-12 13:15:00 +0900 JST
2015/11/12 13:15:00 Next, start at line index 90
```

### Log

Please, look at `trm.log`.

## structure

```
mysql> use trm
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> SHOW TABLES;
+---------------+
| Tables_in_trm |
+---------------+
| conversations |
| tweets        |
| users         |
+---------------+
3 rows in set (0.00 sec)

mysql> DESC users;
+----------+--------------+------+-----+---------+----------------+
| Field    | Type         | Null | Key | Default | Extra          |
+----------+--------------+------+-----+---------+----------------+
| id       | int(11)      | NO   | PRI | NULL    | auto_increment |
| name     | varchar(255) | NO   |     | NULL    |                |
| nickname | varchar(255) | NO   | UNI | NULL    |                |
+----------+--------------+------+-----+---------+----------------+
3 rows in set (0.00 sec)

mysql> DESC tweets;
+------------+--------------+------+-----+---------+----------------+
| Field      | Type         | Null | Key | Default | Extra          |
+------------+--------------+------+-----+---------+----------------+
| id         | int(11)      | NO   | PRI | NULL    | auto_increment |
| twitter_id | varchar(255) | NO   | UNI | NULL    |                |
| success    | int(11)      | NO   |     | NULL    |                |
| user_id    | int(11)      | YES  | MUL | NULL    |                |
| text       | varchar(255) | YES  |     | NULL    |                |
| created_at | timestamp    | YES  |     | NULL    |                |
+------------+--------------+------+-----+---------+----------------+
6 rows in set (0.01 sec)

mysql> DESC conversations;
+-----------------+---------+------+-----+---------+----------------+
| Field           | Type    | Null | Key | Default | Extra          |
+-----------------+---------+------+-----+---------+----------------+
| id              | int(11) | NO   | PRI | NULL    | auto_increment |
| first_tweet_id  | int(11) | NO   | MUL | NULL    |                |
| second_tweet_id | int(11) | NO   | MUL | NULL    |                |
+-----------------+---------+------+-----+---------+----------------+
3 rows in set (0.00 sec)
```

![https://gyazo.com/57d73011d2c9135f8a13f7e5759dbef6](https://i.gyazo.com/57d73011d2c9135f8a13f7e5759dbef6.png)

## notice

### Importing local package

Do not put this repository under `GOPATH` folder because of importing local package.
