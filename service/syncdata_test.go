package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncStockAllFromMongoDBToMySQL(t *testing.T) {
	SyncStockAllFromMongoDBToMySQL()
}

func TestIbuildQuoteDayFromMongoDBToMySQL(t *testing.T) {
	_assert := assert.New(t)
	var code = "sz000001"
	var date = "2020-09-15"
	quote, err := buildQuoteDayFromMongoDBToMySQL(code, date)
	_assert.Nil(err)
	t.Logf("Quote: %v", quote.String())
}

func TestSyncQuoteDayFromMongoDBToMySQL(t *testing.T) {
	_assert := assert.New(t)
	var date = "2020-12-02"
	_, err := SyncQuoteDayFromMongoDBToMySQL(date)
	_assert.Nil(err)
}
