package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInsertQuoteDayManyForMySQL(t *testing.T) {
	Convey("Test Insert QuoteDay Many", t, func() {

		date, err := time.ParseInLocation("2006-01-02", "2020-10-02", time.Local)
		fmt.Println(date)
		So(err, ShouldBeNil)

		var quotes = []*QuoteDay{
			&QuoteDay{
				Code:      "sz000001",
				Open:      10.01,
				Close:     11.01,
				High:      11.01,
				Low:       10.01,
				Volume:    100023123,
				Account:   2314144.12,
				Date:      date,
				DayOfYear: 245,
			},
			&QuoteDay{
				Code:      "sz000002",
				Open:      12.01,
				Close:     12.01,
				High:      12.01,
				Low:       12.01,
				Volume:    100023123,
				Account:   2314144.12,
				Date:      date,
				DayOfYear: 245,
			},
		}

		Convey("Case 1", func() {
			affected, err := InsertQuoteDayManyForMySQL(db.MySQL, quotes)
			So(err, ShouldBeNil)
			So(affected, ShouldEqual, 2)
		})

	})
}
