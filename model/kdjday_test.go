package model

import (
	"testing"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

var kd1 = &KDJDay{
	Code:      "sz000001",
	K:         10.20,
	D:         11.21,
	J:         12.22,
	Date:      time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
	DayOfYear: 345,
}

var kd2 = &KDJDay{
	Code:      "sz000002",
	K:         20.20,
	D:         21.21,
	J:         22.22,
	Date:      time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
	DayOfYear: 345,
}

func TestSelectKDJDayByCodeDateLimit(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	_, err = DeleteKDJDayByCodesDate(tx, []string{kd1.Code}, "2020-12-20")
	_assert.Nil(err)

	var bs = []*KDJDay{kd1}
	_, err = InsertKDJDayMany(tx, bs)
	_assert.Nil(err)
	tx.Commit()

	kdjs, err := SelectKDJDayByCodeDateLimit(db.MySQL, kd1.Code, "2020-12-20", 1)
	_assert.Nil(err)
	_assert.Equal(len(kdjs), 1)
}

func TestInsertKDJDayMany(t *testing.T) {
	_assert := assert.New(t)

	var bs = []*KDJDay{kd1, kd2}

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	affected, err := InsertKDJDayMany(tx, bs)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	tx.Commit()

}

func TestDeleteKDJDayByCodesDate(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)

	_, err = DeleteKDJDayByCodesDate(tx, []string{kd1.Code, kd2.Code}, "2020-12-20")
	_assert.Nil(err)
	tx.Commit()

	tx, err = db.MySQL.Begin()
	_assert.Nil(err)

	var bs = []*KDJDay{kd1, kd2}
	affected, err := InsertKDJDayMany(tx, bs)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = DeleteKDJDayByCodesDate(tx, []string{kd1.Code}, "2020-12-20")
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	affected, err = InsertKDJDayMany(tx, bs)
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = DeleteKDJDayByCodesDate(tx, []string{kd1.Code, kd2.Code}, "2020-12-20")
	_assert.Nil(err)
	_assert.Equal(int64(3), affected)

	tx.Commit()

}
