package model

import (
	"database/sql"
	"testing"

	"github.com/eviltomorrow/aphrodite-calculate/db"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQueryStockListForMongoDB(t *testing.T) {
	Convey("Test Query Stock List For MongoDB", t, func() {
		Convey("Case 1: offset: 0, limit: 20, expect; 20", func() {
			stocks, err := SelectStockListForMongoDB(db.MongoDB, 0, 20)
			So(err, ShouldBeNil)
			So(len(stocks), ShouldEqual, 20)
		})

		Convey("Case 2: offset: 0, limit: 0, expect; 0", func() {
			stocks, err := SelectStockListForMongoDB(db.MongoDB, 0, 0)
			So(err, ShouldBeNil)
			So(len(stocks), ShouldEqual, 0)
		})

		Convey("Case 3: offset: 0, limit: -10, expect; 0", func() {
			stocks, err := SelectStockListForMongoDB(db.MongoDB, 0, -10)
			So(err, ShouldBeNil)
			So(len(stocks), ShouldEqual, 0)
		})

		Convey("Case 4: offset: -10, limit: 20, expect; 0", func() {
			_, err := SelectStockListForMongoDB(db.MongoDB, -10, 20)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestInsertStockManyForMySQL(t *testing.T) {
	Convey("Test Insert Stock Many For MySQL", t, func() {
		var stocks = make([]*Stock, 0, 8)
		stocks = append(stocks, stock1)
		stocks = append(stocks, stock2)
		stocks = append(stocks, stock3)
		stocks = append(stocks, stock4)
		Convey("Case 1: ", func() {
			var insert = make([]*Stock, 0, 8)
			for _, stock := range stocks {
				_, err := SelectStockOneForMySQL(db.MySQL, stock.Code)
				if err == sql.ErrNoRows {
					insert = append(insert, stock)
				}
			}

			affected, err := InsertStockManyForMySQL(db.MySQL, insert)
			So(err, ShouldBeNil)
			So(int64(len(insert)), ShouldEqual, affected)
		})
	})
}

var stock1 = &Stock{
	Code:   "sz000001",
	Name:   "平安银行",
	Source: "sina",
	Valid:  true,
}

var stock2 = &Stock{
	Code:   "sz000002",
	Name:   "万科A",
	Source: "sina",
	Valid:  true,
}
var stock3 = &Stock{
	Code:   "sz000004",
	Name:   "国农科技",
	Source: "sina",
	Valid:  true,
}

var stock4 = &Stock{
	Code:   "sz000005",
	Name:   "世纪星源",
	Source: "sina",
	Valid:  true,
}
