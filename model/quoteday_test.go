package model

import (
	"runtime"
	"testing"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

var date = time.Date(2020, 12, 03, 0, 0, 0, 0, time.Local)

var q1 = &QuoteDay{
	Code:            "sz000001",
	Open:            20.03,
	Close:           21.23,
	High:            21.54,
	Low:             19.33,
	Volume:          521563,
	Account:         20.69 * 521563,
	Date:            date,
	DayOfYear:       date.YearDay(),
	CreateTimestamp: time.Now(),
}

var q2 = &QuoteDay{
	Code:            "sh600365",
	Open:            40.74,
	Close:           43.93,
	High:            44.00,
	Low:             38.63,
	Volume:          1563,
	Account:         41.68 * 1563,
	Date:            date,
	DayOfYear:       date.YearDay(),
	CreateTimestamp: time.Now(),
}

func TestInsertQuoteDayMany(t *testing.T) {
	_assert := assert.New(t)
	tx, err := db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Begin tx error: %v\r\n", err)
	}

	// clear old data
	DeleteQuoteDayByCodesDate(tx, []string{q1.Code, q2.Code}, date.Format("2006-01-02"))

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
	_, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code, q2.Code}, date.Format("2006-01-02"))
	_assert.Nil(err)
	_, err = InsertQuoteDayMany(tx, []*QuoteDay{q1, q2})
	_assert.Nil(err)

	affected, err := DeleteQuoteDayByCodesDate(tx, []string{}, date.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)

	affected, err = DeleteQuoteDayByCodesDate(tx, []string{"no123"}, date.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)

	affected, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code}, date.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	affected, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code}, date.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)

	affected, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code, q2.Code}, date.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	tx.Commit()
}

func TestSelectQuoteDayLatestByCodeDate(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Begin tx error: %v\r\n", err)
	}
	// prepare data
	_, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code, q2.Code}, date.Format("2006-01-02"))
	_assert.Nil(err)
	tx.Commit()

	quotes, err := SelectQuoteDayLatestByCodeDate(db.MySQL, q1.Code, date.Format("2006-01-02"), 30)
	_assert.Nil(err)
	_assert.Equal(0, len(quotes))

	tx, err = db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Begin tx error: %v\r\n", err)
	}
	// prepare data
	_, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code, q2.Code}, date.Format("2006-01-02"))
	_assert.Nil(err)
	_, err = InsertQuoteDayMany(tx, []*QuoteDay{q1})
	_assert.Nil(err)

	tx.Commit()

	quotes, err = SelectQuoteDayLatestByCodeDate(db.MySQL, q1.Code, date.Format("2006-01-02"), 30)
	_assert.Nil(err)
	_assert.Equal(1, len(quotes))
	t.Logf("Open: %f\r\n", quotes[0].Open)

	tx, err = db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Begin tx error: %v\r\n", err)
	}

	for i := 0; i < 100; i++ {
		InsertQuoteDayMany(tx, []*QuoteDay{q1})
	}
	tx.Commit()

	quotes, err = SelectQuoteDayLatestByCodeDate(db.MySQL, q1.Code, date.Format("2006-01-02"), 30)
	_assert.Nil(err)
	_assert.Equal(30, len(quotes))
}

func TestSelectQuoteDayByCodeDate(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Begin tx error: %v\r\n", err)
	}
	// prepare data
	_, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code, q2.Code}, date.Format("2006-01-02"))
	_assert.Nil(err)
	tx.Commit()

	quotes, err := SelectQuoteDayByCodesDate(db.MySQL, []string{q1.Code, q1.Code}, date.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(0, len(quotes))

	tx, err = db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Begin tx error: %v\r\n", err)
	}
	// prepare data
	_, err = DeleteQuoteDayByCodesDate(tx, []string{q1.Code, q2.Code}, date.Format("2006-01-02"))
	_assert.Nil(err)
	_, err = InsertQuoteDayMany(tx, []*QuoteDay{q1, q2})
	_assert.Nil(err)

	tx.Commit()

	quotes, err = SelectQuoteDayByCodesDate(db.MySQL, []string{q1.Code, q2.Code}, date.Format("2006-01-02"))
	_assert.Nil(err)
	_assert.Equal(2, len(quotes))
}

func BenchmarkSelectQuoteDayLatestByCodeDate(b *testing.B) {
	tx, err := db.MySQL.Begin()
	if err != nil {
		b.Fatalf("Begin tx error: %v\r\n", err)
	}

	DeleteQuoteDayByCodesDate(tx, []string{q1.Code}, date.Format("2006-01-02"))
	for i := 0; i < 100; i++ {
		InsertQuoteDayMany(tx, []*QuoteDay{q1})
	}
	tx.Commit()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SelectQuoteDayLatestByCodeDate(db.MySQL, q1.Code, date.Format("2006-01-02"), 30)
	}

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
