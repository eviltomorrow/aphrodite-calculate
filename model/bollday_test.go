package model

import (
	"testing"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

var b1 = &BollDay{
	Code:      "sz000001",
	UP:        10.20,
	MB:        11.21,
	DN:        12.22,
	Date:      time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
	DayOfYear: 345,
}

var b2 = &BollDay{
	Code:      "sz000002",
	UP:        20.20,
	MB:        21.21,
	DN:        22.22,
	Date:      time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
	DayOfYear: 345,
}

func TestInsertBollDayMany(t *testing.T) {
	_assert := assert.New(t)

	var bs = []*BollDay{b1, b2}

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	affected, err := InsertBollDayMany(tx, bs)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	tx.Commit()

}

func TestDeleteBollDayByCodesDate(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	_, err = DeleteBollDayByCodesDate(tx, []string{b1.Code, b2.Code}, "2020-12-20")
	_assert.Nil(err)
	tx.Commit()

	tx, err = db.MySQL.Begin()
	_assert.Nil(err)

	var bs = []*BollDay{b1, b2}
	affected, err := InsertBollDayMany(tx, bs)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = DeleteBollDayByCodesDate(tx, []string{b1.Code}, "2020-12-20")
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	affected, err = InsertBollDayMany(tx, bs)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = DeleteBollDayByCodesDate(tx, []string{b1.Code, b2.Code}, "2020-12-20")
	_assert.Nil(err)
	_assert.Equal(int64(3), affected)

	tx.Commit()

}
