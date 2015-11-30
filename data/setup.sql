create table dev_tweet_ids (
  post_id varchar(32) not null,
  reply_id varchar(32) not null,
  label_01 varchar(4) not null,
  label_02 varchar(4) not null,
  label_03 varchar(4) not null,
  label_04 varchar(4) not null,
  label_05 varchar(4) not null,
  label_06 varchar(4) not null,
  label_07 varchar(4) not null,
  label_08 varchar(4) not null,
  label_09 varchar(4) not null,
  label_10 varchar(4) not null
);

load data local infile '~/go/src/github.com/sosuke-k/twitter-responding-machine/data/dev.txt' into table dev_tweet_ids;
