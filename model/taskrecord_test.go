package model

import (
	"fmt"
	"testing"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

var t1 = &TaskRecord{
	Name:      "同步日数据",
	Code:      "sz000001",
	Date:      "2020-12-10",
	Completed: false,
	Msg:       "",
}

var t2 = &TaskRecord{
	Name:      "同步日数据",
	Code:      "sz000002",
	Date:      "2020-12-10",
	Completed: false,
	Msg:       "",
}

var t3 = &TaskRecord{
	Name:      "同步日数据",
	Code:      "sh000536",
	Date:      "2020-12-10",
	Completed: false,
	Msg:       "",
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
			Name: "",
			Date: "2020-12-10",
		}

		t.Code = fmt.Sprintf("sh6000%d", i)
		records = append(records, t)
	}

	affected, err = InsertTaskRecordMany(tx, records)
	_assert.Nil(err)
	_assert.Equal(int64(30), affected)

}
