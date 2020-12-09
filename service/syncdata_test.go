package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsyncStockFromMongoDBToMySQL(t *testing.T) {
	_assert := assert.New(t)

	var offset int64 = 0
	var limit int64 = 30
	var lastID string
	affected, lastID, err := syncStockFromMongoDBToMySQL(offset, limit, lastID)
	_assert.Nil(err)
	_assert.NotEqual("", lastID)
	t.Logf("Affetced: %v\r\n", affected)

	offset = 0
	limit = 0
	affected, lastID, err = syncStockFromMongoDBToMySQL(offset, limit, lastID)
	_assert.Nil(err)
	_assert.Equal("", lastID)
	_assert.Equal(int64(0), affected)
}

func TestSyncStockAllFromMongoDBToMySQL(t *testing.T) {
	// _assert := assert.New(t)
	SyncStockAllFromMongoDBToMySQL()
}

func BenchmarkIsyncStockFromMongoDBToMySQL(b *testing.B) {
	var offset int64 = 0
	var limit int64 = 30
	var lastID string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		syncStockFromMongoDBToMySQL(offset, limit, lastID)
	}
}

func BenchmarkSyncStockAllFromMongoDBToMySQL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SyncStockAllFromMongoDBToMySQL()
	}
}
