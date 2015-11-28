package twitter_test

import "testing"
import (
	"../twitter"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTweet(t *testing.T) {

	Convey("fetch without id", t, func() {
		tweet := twitter.Tweet{}
		err := tweet.Fetch()
		So(err, ShouldNotEqual, nil)

		//check my custom options
		So(tweet.Success, ShouldEqual, 0)
		So(err.(*twitter.Error).Op, ShouldEqual, twitter.Op.Query)
	})

	Convey("fetch by tweet id\n", t, func() {
		// https://twitter.com/olha_drm/status/418033807850496002
		tweet := twitter.Tweet{ItemID: "418033807850496002"}
		err := tweet.Fetch()
		So(err, ShouldEqual, nil)

		//check tweet and reply
		So(tweet.Success, ShouldEqual, 1)
		So(tweet.Text, ShouldEqual, "よるほ　あけましておめでとうございますほー")
		So(tweet.ScreenName, ShouldEqual, "olha_drm")
		So(tweet.Name, ShouldEqual, "織羽")
		So(len(tweet.Replies), ShouldEqual, 1)
		So(tweet.Replies[0].Text, ShouldEqual, "@olha_drm あけおめっ！ことよろー！！")
		So(tweet.Replies[0].ScreenName, ShouldEqual, "reprohonmono")
		// So(tweet.Replies[0].Name, ShouldEqual, "環境科学コース")
		So(tweet.Replies[0].ReplyTo, ShouldEqual, "418033807850496002")
	})

	Convey("cannot fetch because of authorization\n", t, func() {
		tweet := twitter.Tweet{ItemID: "418033823511629836"}
		err := tweet.Fetch()
		So(err, ShouldNotEqual, nil)

		//check my custom options
		So(tweet.Success, ShouldEqual, 0)
		So(err.(*twitter.Error).Op, ShouldEqual, twitter.Op.Authorization)
	})

	Convey("cannot fetch because of not exist\n", t, func() {
		tweet := twitter.Tweet{ItemID: "418033823511629837"}
		err := tweet.Fetch()
		So(err, ShouldNotEqual, nil)

		//check my custom options
		So(tweet.Success, ShouldEqual, 0)
		So(err.(*twitter.Error).Op, ShouldEqual, twitter.Op.NotExisting)
	})

}
