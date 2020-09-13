package model

import (
	"fmt"
	"testing"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

func TestQueryStockListForMongoDB(t *testing.T) {
	_assert := assert.New(t)

	stocks, err := QueryStockListForMongoDB(db.MongoDB, 2, 1)
	_assert.Nil(err)
	t.Log(stocks)
}

func TestInsertStockManyForMySQL(t *testing.T) {
	_assert := assert.New(t)
	var stocks = make([]*Stock, 0, 3)
	for i := 0; i < 3; i++ {
		stocks = append(stocks, &Stock{
			Source: "Sina",
			Code:   fmt.Sprintf("sz00000%d", i),
			Name:   "瑞幸咖啡",
		})
	}

	affected, err := InsertStockManyForMySQL(db.MySQL, stocks)
	_assert.Nil(err)
	_assert.Equal(int64(3), affected)
}

func TestQueryStockListForMySQL(t *testing.T) {
	_assert := assert.New(t)
	stocks, err := QueryStockListForMySQL(db.MySQL, 1, 10)
	_assert.Nil(err)
	for _, stock := range stocks {
		t.Logf("%s\r\n", stock.String())
	}
}
