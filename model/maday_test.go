package model

import (
	"testing"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

var ma1 = &MADay{
	Code:      "sz000001",
	M5:        10.20,
	M10:       11.21,
	M20:       12.22,
	M30:       13.23,
	Date:      time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
	DayOfYear: 345,
}

var ma2 = &MADay{
	Code:      "sz000002",
	M5:        20.20,
	M10:       21.21,
	M20:       22.22,
	M30:       23.23,
	Date:      time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
	DayOfYear: 345,
}

func TestInsertMADayMany(t *testing.T) {
	_assert := assert.New(t)

	var mas = []*MADay{ma1, ma2}

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	affected, err := InsertMADayMany(tx, mas)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	tx.Commit()

}

func TestDeleteMADayByCodesDate(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	_, err = DeleteMADayByCodesDate(tx, []string{ma1.Code, ma2.Code}, "2020-12-20")
	_assert.Nil(err)
	tx.Commit()

	tx, err = db.MySQL.Begin()
	_assert.Nil(err)

	var mas = []*MADay{ma1, ma2}
	affected, err := InsertMADayMany(tx, mas)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = DeleteMADayByCodesDate(tx, []string{ma1.Code}, "2020-12-20")
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	affected, err = InsertMADayMany(tx, mas)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = DeleteMADayByCodesDate(tx, []string{ma1.Code, ma2.Code}, "2020-12-20")
	_assert.Nil(err)
	_assert.Equal(int64(3), affected)

	tx.Commit()

}
