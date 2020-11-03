package model

import (
	"testing"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	. "github.com/smartystreets/goconvey/convey"
)

func TestQueryQuoteOne(t *testing.T) {

	Convey("Test QueryQuoteOne", t, func() {
		Convey("Collect case:", func() {
			var where = map[string]interface{}{
				"code": "sz000001",
				"date": "2020-09-15",
			}

			quote, err := QueryQuoteBaseOne(db.MongoDB, where)
			So(err, ShouldBeNil)
			So(quote, ShouldNotBeNil)
		})
	})

}
