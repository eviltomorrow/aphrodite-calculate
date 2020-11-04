package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSelectQuoteDayByCodesDateForMySQL(t *testing.T) {
	Convey("Test Delete QuoteDay Many", t, func() {
		Convey("Case1", func() {
			quotes, err := SelectQuoteDayByCodeDate(db.MySQL, []string{"sz000002"}, "2020-10-02")
			So(err, ShouldBeNil)
			So(len(quotes), ShouldEqual, 4)
			for _, quote := range quotes {
				t.Logf("quote: %v\r\n", quote.String())
			}
		})
	})
}

func TestSelectQuoteDayLatestByCodeDate(t *testing.T) {
	Convey("Test Delete QuoteDay Many", t, func() {
		Convey("Case1", func() {
			quotes, err := SelectQuoteDayLatestByCodeDate(db.MySQL, "sz000002", "2020-10-02", 10)
			So(err, ShouldBeNil)
			So(len(quotes), ShouldEqual, 4)
			for _, quote := range quotes {
				t.Logf("quote: %v\r\n", quote.String())
			}
		})
	})
}

func TestDeleteQuoteDayByCodeDateForMySQL(t *testing.T) {
	Convey("Test Delete QuoteDay Many", t, func() {
		Convey("Case1", func() {
			affected, err := DeleteQuoteDayByCodeDate(db.MySQL, "sz000002", "2020-10-02")
			So(err, ShouldBeNil)
			So(affected, ShouldEqual, 1)
		})
	})
}

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
			affected, err := InsertQuoteDayMany(db.MySQL, quotes)
			So(err, ShouldBeNil)
			So(affected, ShouldEqual, 2)
		})

	})
}
