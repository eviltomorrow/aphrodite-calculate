package model

import (
	"testing"
	"time"

	"github.com/eviltomorrow/aphrodite-base/ztime"
	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

var beginDate = time.Date(2020, 11, 30, 0, 0, 0, 0, time.Local)
var endDate = time.Date(2020, 12, 04, 0, 0, 0, 0, time.Local)

var qw1 = &QuoteWeek{
	Code:            "sz000001",
	Open:            20.30,
	Close:           25.23,
	High:            27.54,
	Low:             16.33,
	Volume:          12521563,
	Account:         20.69 * 12521563,
	DateBegin:       beginDate,
	DateEnd:         endDate,
	WeekOfYear:      ztime.YearWeek(endDate),
	CreateTimestamp: time.Now(),
}

var qw2 = &QuoteWeek{
	Code:            "sz000002",
	Open:            30.30,
	Close:           35.23,
	High:            47.54,
	Low:             26.33,
	Volume:          42521563,
	Account:         30.69 * 42521563,
	DateBegin:       beginDate,
	DateEnd:         endDate,
	WeekOfYear:      ztime.YearWeek(endDate),
	CreateTimestamp: time.Now(),
}

func TestDeleteQuoteWeekByCodesDate(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Begin error: %v\r\n", err)
	}
	_, err = DeleteQuoteWeekByCodesDate(tx, []string{qw1.Code, qw2.Code}, endDate.Format("2006-01-02"))
	_assert.Nil(err)
	tx.Commit()

	tx, err = db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Begin error: %v\r\n", err)
	}
	affected, err := DeleteQuoteWeekByCodesDate(tx, []string{qw1.Code, qw2.Code}, endDate.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)
	tx.Commit()
}

func TestInsertQuoteWeekMany(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Begin error: %v\r\n", err)
	}
	_, err = DeleteQuoteWeekByCodesDate(tx, []string{qw1.Code, qw2.Code}, endDate.Format("2006-01-02"))
	_assert.Nil(err)

	affected, err := InsertQuoteWeekMany(tx, []*QuoteWeek{qw1, qw2})
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = InsertQuoteWeekMany(tx, []*QuoteWeek{})
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)

	var quotes = make([]*QuoteWeek, 0, 20)
	for i := 0; i < 20; i++ {
		quotes = append(quotes, qw1)
	}

	affected, err = InsertQuoteWeekMany(tx, quotes)
	_assert.Nil(err)
	_assert.Equal(int64(20), affected)

	tx.Commit()
}

func TestSelectQuoteWeekByCodeDate(t *testing.T) {
	_assert := assert.New(t)

	tx, _ := db.MySQL.Begin()
	_, err := DeleteQuoteWeekByCodesDate(tx, []string{qw1.Code, qw2.Code}, endDate.Format("2006-01-02"))
	_assert.Nil(err)
	tx.Commit()

	tx, _ = db.MySQL.Begin()
	affected, err := InsertQuoteWeekMany(tx, []*QuoteWeek{qw1, qw2})
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)
	tx.Commit()

	quotes, err := SelectQuoteWeekByCodeDate(db.MySQL, []string{qw1.Code, qw2.Code}, endDate.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(2, len(quotes))
}
