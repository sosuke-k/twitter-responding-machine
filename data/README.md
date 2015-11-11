# Generating Tweet Database

I used [GORM](https://github.com/jinzhu/gorm), MySQL and [twittergo](https://github.com/kurrik/twittergo) and referred to [twittergo-examples](https://github.com/kurrik/twittergo-examples).

## Structure

![https://gyazo.com/57d73011d2c9135f8a13f7e5759dbef6](https://i.gyazo.com/57d73011d2c9135f8a13f7e5759dbef6.png)

## preliminary

### Tweet ID Dataset

```
wget https://github.com/mynlp/stc/raw/master/taskdata/twitter_id_str_data.txt.bz2
bzip2 -d twitter_id_str_data.txt.bz2
```

## Twitter API

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

## using

```
go run main.go
```
