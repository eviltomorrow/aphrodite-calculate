package model

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/eviltomorrow/aphrodite-calculate/db"
)

var date1 = time.Date(2020, 12, 03, 0, 0, 0, 0, time.Local)
var date2 = time.Date(2020, 12, 04, 0, 0, 0, 0, time.Local)

var q1 = &QuoteDay{
	Code:            "sz000001",
	Open:            20.03,
	Close:           21.23,
	High:            21.54,
	Low:             19.33,
	YesterdayClosed: 16.23,
	Volume:          521563,
	Account:         20.69 * 521563,
	Date:            date1,
	DayOfYear:       date1.YearDay(),
	CreateTimestamp: time.Now(),
}

var q2 = &QuoteDay{
	Code:            "sh600365",
	Open:            50.74,
	Close:           42.93,
	High:            43.00,
	Low:             38.63,
	YesterdayClosed: 16.23,
	Volume:          14563,
	Account:         41.68 * 14563,
	Date:            date1,
	DayOfYear:       date1.YearDay(),
	CreateTimestamp: time.Now(),
}

var q3 = &QuoteDay{
	Code:            "sh600365",
	Open:            40.74,
	Close:           43.93,
	High:            44.00,
	Low:             38.63,
	YesterdayClosed: 16.23,
	Volume:          1563,
	Account:         41.68 * 1563,
	Date:            date2,
	DayOfYear:       date2.YearDay(),
	CreateTimestamp: time.Now(),
}

func TestInsertQuoteDayMany(t *testing.T) {
	_assert := assert.New(t)
	tx, err := db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Begin tx error: %v\r\n", err)
	}

	// clear old data
	DeleteQuoteDayByCodesDate(tx, []string{q1.Code, q2.Code}, date1.Format("2006-01-02"))

	// right
	affected, err := InsertQuoteDayMany(tx, []*QuoteDay{q1, q2})
	_assert.Nil(err)
	_assert.Equal(int64(2), affected)

	affected, err = InsertQuoteDayMany(tx, []*QuoteDay{})
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)

	tx.Commit()
}

func TestDeleteQuoteDayByCodesDate(t *testing.T) {
	_assert := assert.New(t)
	tx, err := db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Begin tx error: %v\r\n", err)
	}

	// prepare data
	_, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code, q2.Code}, date1.Format("2006-01-02"))
	_assert.Nil(err)
	_, err = InsertQuoteDayMany(tx, []*QuoteDay{q1, q2})
	_assert.Nil(err)

	affected, err := DeleteQuoteDayByCodesDate(tx, []string{}, date1.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)

	affected, err = DeleteQuoteDayByCodesDate(tx, []string{"no123"}, date1.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)

	affected, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code}, date1.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	affected, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code}, date1.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)

	affected, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code, q2.Code}, date1.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	tx.Commit()
}

func TestSelectQuoteDayByCodeDate(t *testing.T) {
	_assert := assert.New(t)
	tx, err := db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Begin tx error: %v\r\n", err)
	}

	// prepare data
	_, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code, q2.Code}, date1.Format("2006-01-02"))
	_assert.Nil(err)
	_, err = InsertQuoteDayMany(tx, []*QuoteDay{q1, q2})
	_assert.Nil(err)

	_, err = DeleteQuoteDayByCodesDate(tx, []string{q3.Code}, date2.Format("2006-01-02"))
	_assert.Nil(err)
	_, err = InsertQuoteDayMany(tx, []*QuoteDay{q3})
	_assert.Nil(err)
	tx.Commit()

	quotes, err := SelectQuoteDayByCodeDate(db.MySQL, q3.Code, date1.Format("2006-01-02"), date2.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(2, len(quotes))

	quotes, err = SelectQuoteDayByCodeDate(db.MySQL, q1.Code, date1.Format("2006-01-02"), date2.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(1, len(quotes))

	quotes, err = SelectQuoteDayByCodeDate(db.MySQL, q1.Code, time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(0, len(quotes))

	_assert.Nil(nil)
}

func BenchmarkInsertQuoteDayMany(b *testing.B) {
	tx, err := db.MySQL.Begin()
	if err != nil {
		b.Fatalf("Error: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// DeleteQuoteDayByCodesDate(tx, []string{q1.Code, q2.Code}, date.Format("2006-01-02"))
		InsertQuoteDayMany(tx, []*QuoteDay{q1, q2})
	}
	tx.Commit()
}

func BenchmarkParallelInsertQuoteDayMany(b *testing.B) {
	b.SetParallelism(runtime.NumCPU())
	tx, err := db.MySQL.Begin()
	if err != nil {
		b.Fatalf("Error: %v", err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			InsertQuoteDayMany(tx, []*QuoteDay{q1, q2})
		}
	})
	tx.Commit()
}
