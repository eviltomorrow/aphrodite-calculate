package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSyncStockAllFromMongoDBToMySQL(t *testing.T) {
	SyncStockAllFromMongoDBToMySQL()
}

func TestSyncQuoteDayFromMongoDBToMySQL(t *testing.T) {
	_assert := assert.New(t)
	var date = "2020-12-02"
	_, err := SyncQuoteDayFromMongoDBToMySQL(date)
	_assert.Nil(err)
}

func TestSyncQuoteWeekFromMongoDBToMySQL(t *testing.T) {
	_assert := assert.New(t)
	// var date = "2020-11-23"
	// _, err := SyncQuoteDayFromMongoDBToMySQL(date)
	// _assert.Nil(err)

	// date = "2020-11-24"
	// _, err = SyncQuoteDayFromMongoDBToMySQL(date)
	// _assert.Nil(err)

	// date = "2020-11-25"
	// _, err = SyncQuoteDayFromMongoDBToMySQL(date)
	// _assert.Nil(err)

	// date = "2020-11-26"
	// _, err = SyncQuoteDayFromMongoDBToMySQL(date)
	// _assert.Nil(err)

	// date = "2020-11-27"
	// _, err = SyncQuoteDayFromMongoDBToMySQL(date)
	// _assert.Nil(err)

	var date = "2020-11-27"
	_, err := SyncQuoteWeekFromMongoDBToMySQL(date)
	_assert.Nil(err)
}

func TestIbuildQuoteWeekFromQuoteDay(t *testing.T) {
	_assert := assert.New(t)

	begin, err := time.ParseInLocation("2006-01-02", "2020-11-23", time.Local)
	_assert.Nil(err)

	end, err := time.ParseInLocation("2006-01-02", "2020-11-27", time.Local)
	_assert.Nil(err)

	data, err := buildQuoteWeekFromQuoteDay("sz300999", begin, end)
	_assert.Nil(err)
	t.Logf("Data: %s\r\n", data)
}
