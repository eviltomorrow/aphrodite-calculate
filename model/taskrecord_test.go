package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eviltomorrow/aphrodite-calculate/db"
)

var t1 = &TaskRecord{
	Method:    "SYNC_QUOTEDAY",
	Date:      "2020-12-10",
	Completed: false,
}

var t2 = &TaskRecord{
	Method:    "SYNC_QUOTEDAY",
	Date:      "2020-12-10",
	Completed: false,
}

var t3 = &TaskRecord{
	Method:    "SYNC_QUOTEWEEK",
	Date:      "2020-12-10",
	Completed: false,
}

func TestSelectTaskRecordManyByDate(t *testing.T) {
	_assert := assert.New(t)
	records, err := SelectTaskRecordManyByDate(db.MySQL, "2020-12-02")
	_assert.Nil(err)
	_assert.Equal(2, len(records))
}

func TestInsertTaskRecordMany(t *testing.T) {
	_assert := assert.New(t)

	tx, err := db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Error: %v\r\n", err)
	}

	affected, err := InsertTaskRecordMany(tx, []*TaskRecord{t1, t2, t3})
	_assert.Nil(err)
	_assert.Equal(int64(3), affected)

	affected, err = InsertTaskRecordMany(tx, []*TaskRecord{})
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)

	var records = make([]*TaskRecord, 0, 30)
	for i := 10; i < 40; i++ {
		var t = &TaskRecord{
			Method: "SYNC_QUOTEDAY",
			Date:   "2020-12-10",
		}
		records = append(records, t)
	}

	affected, err = InsertTaskRecordMany(tx, records)
	_assert.Nil(err)
	_assert.Equal(int64(30), affected)
	tx.Commit()
}

func TestUpdateTaskRecordCompleted(t *testing.T) {
	_assert := assert.New(t)

	var date = "2020-12-10"
	records, err := SelectTaskRecordManyByDate(db.MySQL, date)
	_assert.Nil(err)

	tx, err := db.MySQL.Begin()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	for _, record := range records {
		UpdateTaskRecordCompleted(tx, []int64{record.ID})
	}

	tx.Commit()
}

func BenchmarkSelectTaskRecordManyByDate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SelectTaskRecordManyByDate(db.MySQL, "2020-12-02")
	}
}

func BenchmarkInsertTaskRecordMany(b *testing.B) {
	var records = []*TaskRecord{
		t1, t2, t3,
	}

	tx, err := db.MySQL.Begin()
	if err != nil {
		b.Fatalf("Error: %v\r\n", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		InsertTaskRecordMany(tx, records)
	}
	tx.Commit()
}
