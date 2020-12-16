package model

import (
	"testing"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

func TestInsertTaskRecordMany(t *testing.T) {
	_assert := assert.New(t)

	var record = &TaskRecord{
		Method:     "SYNC_QUOTEDAY",
		Date:       "2020-12-03",
		Priority:   0,
		Completed:  false,
		NumOfTimes: 0,
	}

	tx, err := db.MySQL.Begin()
	_assert.Nil(err)
	affected, err := InsertTaskRecordMany(tx, []*TaskRecord{record})
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	tx.Commit()
}
