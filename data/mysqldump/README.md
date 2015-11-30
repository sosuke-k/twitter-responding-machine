# Dump Dataset for STC

## Setup

```
$ mysql -u root

mysql> CREATE DATABASE <database_name> DEFAULT CHARACTER SET utf8;
Query OK, 1 row affected (0.00 sec)

mysql> quit;
Bye

$ mysql -u root <database_name> < stc.dump
```

## Check

```
$ mysql -u root

mysql> use <database_name>;
Database changed

mysql> show tables;
+---------------+
| Tables_in_stc |
+---------------+
| stc_tweet_ids |
| stc_tweets    |
+---------------+
2 rows in set (0.00 sec)

mysql> select count(*) from stc_tweet_ids;
+----------+
| count(*) |
+----------+
|   500000 |
+----------+
1 row in set (0.17 sec)

mysql> select count(*) from stc_tweets;
+----------+
| count(*) |
+----------+
|  1000000 |
+----------+
1 row in set (0.26 sec)
```

## Architecture

### stc_tweet_ids

```
mysql> desc stc_tweet_ids;
+----------+-------------+------+-----+---------+-------+
| Field    | Type        | Null | Key | Default | Extra |
+----------+-------------+------+-----+---------+-------+
| post_id  | varchar(32) | NO   |     | NULL    |       |
| reply_id | varchar(32) | NO   |     | NULL    |       |
+----------+-------------+------+-----+---------+-------+
2 rows in set (0.01 sec)
```

`(post_id, reply_id)` inserted from [mynlp/stc/taskdata](https://github.com/mynlp/stc/tree/master/taskdata) with no index.
Please if you use this table, do that on your own.


### stc_tweets

```
mysql> desc stc_tweets;
+-------------+--------------+------+-----+---------+----------------+
| Field       | Type         | Null | Key | Default | Extra          |
+-------------+--------------+------+-----+---------+----------------+
| id          | int(11)      | NO   | PRI | NULL    | auto_increment |
| success     | int(11)      | NO   |     | NULL    |                |
| item_id     | varchar(255) | NO   | MUL | NULL    |                |
| screen_name | varchar(255) | YES  | MUL | NULL    |                |
| name        | varchar(255) | YES  |     | NULL    |                |
| time        | varchar(255) | YES  |     | NULL    |                |
| text        | varchar(255) | YES  |     | NULL    |                |
+-------------+--------------+------+-----+---------+----------------+
7 rows in set (0.01 sec)
```

I add index to only `item_id`, `screen_name` columns.
Please optimize data size of each columns on your own.

| Field       | Content                                   |
:-------------|:------------------------------------------|
| id          | surrogate key not related to twitter      |
| success     | success flag i.e. 1, -1, and -2           |
| item_id     | tweet data item id                        |
| screen_name | tweet data screen name                    |
| name        | tweet data name                           |
| time        | tweet data time (raw string of unix time) |
| text        | tweet text                                |


#### success flag

| flag | meaning                                                        |
|:----:|:---------------------------------------------------------------|
| 1    | succeeded                                                      |
| -1   | redirected to URL which not including `item_id`                |
| -2   | fetched HTML which including "Sorry, that page doesn’t exist!" |
| 0    | anything else                                                  |

but no 0 flag is here.

## Valid Data (Taskdata)

### Tweet

There are 921314 valid tweets.

```
mysql> select success, count(*) from stc_tweets group by success;
+---------+----------+
| success | count(*) |
+---------+----------+
|      -2 |    25953 |
|      -1 |    52733 |
|       1 |   921314 |
+---------+----------+
3 rows in set (0.77 sec)
```

### Conversation

There are 427307 valid conversations.

```
mysql> select count(*) from stc_tweets as t1 inner join stc_tweet_ids as ids on t1.item_id = ids.post_id inner join stc_tweets as t2 on ids.reply_id = t2.item_id where t1.success = 1 and t2.success = 1;
+----------+
| count(*) |
+----------+
|   427307 |
+----------+
1 row in set (5.70 sec)
```

## Valid Data (Devset)

### Post Tweet

There are 182/200 valid post tweets.

### Conversation

There are 1651/1959 valid conversations.

* pc = あるポストに対してベースラインシステムが出力したツイートID１０個のうち取得できた数
* pn = `pc` だけ取得できたポストの数

|pc|pn|
|:-:|:-:|
|  6 |        3 |
|  7 |        4 |
|  8 |       29 |
|  9 |       77 |
| 10 |       68 |

1651 = 10 * 68 + 9 * 77 + 8 * 29 + 7 * 4 + 6 * 3
