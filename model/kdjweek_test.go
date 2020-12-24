package model

import (
	"testing"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

var kw1 = &KDJWeek{
	Code:       "sz000001",
	K:          10.20,
	D:          11.21,
	J:          12.22,
	Date:       time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
	WeekOfYear: 50,
}

var kw2 = &KDJWeek{
	Code:       "sz000002",
	K:          20.20,
	D:          21.21,
	J:          22.22,
	Date:       time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
	WeekOfYear: 50,
}

func TestSelectKDJWeekByCodeDateLimit(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	_, err = DeleteKDJWeekByCodesDate(tx, []string{kw1.Code}, "2020-12-20")
	_assert.Nil(err)

	var bs = []*KDJWeek{kw1}
	_, err = InsertKDJWeekMany(tx, bs)
	_assert.Nil(err)
	tx.Commit()

	kdjs, err := SelectKDJWeekByCodeDateLimit(db.MySQL, kw1.Code, "2020-12-20", 1)
	_assert.Nil(err)
	_assert.Equal(len(kdjs), 1)
}

func TestInsertKDJWeekMany(t *testing.T) {
	_assert := assert.New(t)

	var bs = []*KDJWeek{kw1, kw2}

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	affected, err := InsertKDJWeekMany(tx, bs)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	tx.Commit()

}

func TestDeleteKDJWeekByCodesDate(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	_, err = DeleteKDJWeekByCodesDate(tx, []string{kw1.Code, kw2.Code}, "2020-12-20")
	_assert.Nil(err)
	tx.Commit()

	tx, err = db.MySQL.Begin()
	_assert.Nil(err)

	var bs = []*KDJWeek{kw1, kw2}
	affected, err := InsertKDJWeekMany(tx, bs)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = DeleteKDJWeekByCodesDate(tx, []string{kw1.Code}, "2020-12-20")
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	affected, err = InsertKDJWeekMany(tx, bs)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = DeleteKDJWeekByCodesDate(tx, []string{kw1.Code, kw2.Code}, "2020-12-20")
	_assert.Nil(err)
	_assert.Equal(int64(3), affected)

	tx.Commit()

}
