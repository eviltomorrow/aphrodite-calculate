package model

import (
	"testing"

	"github.com/eviltomorrow/aphrodite-calculate/db"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInsertStockManyForMySQL(t *testing.T) {
	Convey("Test Insert Stock Many", t, func() {
		var stocks = []*Stock{
			&Stock{
				Code:            "sz000001",
				Name:            "平安银行",
				Source:          "sina",
				CreateTimestamp: 0,
			},
			&Stock{
				Code:            "sz000002",
				Name:            "平安证券",
				Source:          "sina",
				CreateTimestamp: 0,
			},
		}

		Convey("Case 1", func() {
			affected, err := InsertStockManyForMySQL(db.MySQL, stocks)
			So(err, ShouldBeNil)
			So(affected, ShouldEqual, 2)
		})

		Convey("Case 2", func() {
			affected, err := InsertStockManyForMySQL(db.MySQL, stocks)
			So(err, ShouldNotBeNil)
			So(affected, ShouldEqual, 0)
		})
	})
}

func TestUpdateStockByCodeForMySQL(t *testing.T) {
	Convey("Test Update Stock By Code For MySQL", t, func() {
		var stock = &Stock{
			Code:   "sz000001",
			Name:   "平安资管",
			Source: "sina",
		}
		Convey("Case 1", func() {
			affected, err := UpdateStockByCodeForMySQL(db.MySQL, "sz000001", stock)
			So(err, ShouldBeNil)
			So(affected, ShouldEqual, 1)
		})
	})
}

func TestQueryStockListForMongoDB(t *testing.T) {
	Convey("Test Query Stock List For MongoDB", t, func() {
		Convey("Case 1", func() {
			stocks, err := QueryStockListForMongoDB(db.MongoDB, 0, 20)
			So(err, ShouldBeNil)
			So(len(stocks), ShouldEqual, 20)
		})

		Convey("Case 2", func() {
			stocks, err := QueryStockListForMongoDB(db.MongoDB, 0, 0)
			So(err, ShouldBeNil)
			So(len(stocks), ShouldEqual, 0)
		})
	})
}
