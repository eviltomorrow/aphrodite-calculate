package model

import (
	"testing"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

var mw1 = &MAWeek{
	Code:       "sz000001",
	M5:         10.20,
	M10:        11.21,
	M20:        12.22,
	M30:        13.23,
	Date:       time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
	WeekOfYear: 50,
}

var mw2 = &MAWeek{
	Code:       "sz000002",
	M5:         20.20,
	M10:        21.21,
	M20:        22.22,
	M30:        23.23,
	Date:       time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
	WeekOfYear: 50,
}

func TestInsertMAWeekMany(t *testing.T) {
	_assert := assert.New(t)

	var mws = []*MAWeek{mw1, mw2}

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	affected, err := InsertMAWeekMany(tx, mws)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	tx.Commit()

}

func TestDeleteMAWeekByCodesDate(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	_, err = DeleteMAWeekByCodesDate(tx, []string{mw1.Code, mw2.Code}, "2020-12-20")
	_assert.Nil(err)
	tx.Commit()

	tx, err = db.MySQL.Begin()
	_assert.Nil(err)

	var mws = []*MAWeek{mw1, mw2}
	affected, err := InsertMAWeekMany(tx, mws)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = DeleteMAWeekByCodesDate(tx, []string{mw1.Code}, "2020-12-20")
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	affected, err = InsertMAWeekMany(tx, mws)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = DeleteMAWeekByCodesDate(tx, []string{mw1.Code, mw2.Code}, "2020-12-20")
	_assert.Nil(err)
	_assert.Equal(int64(3), affected)

	tx.Commit()

}
