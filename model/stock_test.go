package model

import (
	"testing"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

var s1 = &Stock{
	Code:            "sz000001",
	Name:            "平安银行",
	Source:          "sina",
	Valid:           true,
	CreateTimestamp: Time(time.Now()),
}

var s2 = &Stock{
	Code:            "sz000002",
	Name:            "平安证券",
	Source:          "sina",
	Valid:           true,
	CreateTimestamp: Time(time.Now()),
}

var s3 = &Stock{
	Code:            "sh600235",
	Name:            "中国软件",
	Source:          "sina",
	Valid:           true,
	CreateTimestamp: Time(time.Now()),
}

func TestInsertStockManyForMySQL(t *testing.T) {
	// _assert := assert.New(t)

}

func TestSelectStockListForMongoDB(t *testing.T) {
	_assert := assert.New(t)

	var offset int64 = -1
	var limit int64 = -1
	stocks, err := SelectStockManyForMongoDB(db.MongoDB, offset, limit, "")
	_assert.Nil(err)
	_assert.Equal(0, len(stocks))

	offset = 0
	limit = 30
	stocks, err = SelectStockManyForMongoDB(db.MongoDB, offset, limit, "")
	_assert.Nil(err)
	_assert.Equal(30, len(stocks))

	offset = 4200
	limit = 30
	stocks, err = SelectStockManyForMongoDB(db.MongoDB, offset, limit, "")
	_assert.Nil(err)
	_assert.Equal(0, len(stocks))

	offset = 1
	limit = 30
	stocks, err = SelectStockManyForMongoDB(db.MongoDB, offset, limit, "5f4cbd3c37ab7a704504e6f0")
	_assert.Nil(err)
	_assert.Equal(30, len(stocks))

	var count int
	var lastID string
	for {
		stocks, err = SelectStockManyForMongoDB(db.MongoDB, 0, 30, lastID)
		if len(stocks) == 0 || len(stocks) < 30 {
			count += len(stocks)
			break
		}
		lastID = stocks[len(stocks)-1].ObjectID
		count += len(stocks)
	}

	t.Logf("Count: %v\r\n", count)
}

func BenchmarkSelectStockListForMongoDB(b *testing.B) {
	var offset int64 = 0
	var limit int64 = 100
	var objectID string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stocks, _ := SelectStockManyForMongoDB(db.MongoDB, offset, limit, objectID)
		if len(stocks) != 0 {
			objectID = stocks[len(stocks)-1].ObjectID
		}
		offset += limit
	}
}
