package model

import (
	"testing"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

var bw1 = &BollWeek{
	Code:       "sz000001",
	UP:         10.20,
	MB:         11.21,
	DN:         12.22,
	Date:       time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
	WeekOfYear: 50,
}

var bw2 = &BollWeek{
	Code:       "sz000002",
	UP:         20.20,
	MB:         21.21,
	DN:         22.22,
	Date:       time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
	WeekOfYear: 50,
}

func TestInsertBollWeekMany(t *testing.T) {
	_assert := assert.New(t)

	var bs = []*BollWeek{bw1, bw2}

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	affected, err := InsertBollWeekMany(tx, bs)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	tx.Commit()

}

func TestDeleteBollWeekByCodesDate(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	_, err = DeleteBollWeekByCodesDate(tx, []string{b1.Code, b2.Code}, "2020-12-20")
	_assert.Nil(err)
	tx.Commit()

	tx, err = db.MySQL.Begin()
	_assert.Nil(err)

	var bs = []*BollWeek{bw1, bw2}
	affected, err := InsertBollWeekMany(tx, bs)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = DeleteBollWeekByCodesDate(tx, []string{b1.Code}, "2020-12-20")
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	affected, err = InsertBollWeekMany(tx, bs)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = DeleteBollWeekByCodesDate(tx, []string{b1.Code, b2.Code}, "2020-12-20")
	_assert.Nil(err)
	_assert.Equal(int64(3), affected)

	tx.Commit()

}
