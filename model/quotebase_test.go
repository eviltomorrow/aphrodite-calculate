package model

import (
	"testing"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

func TestQueryQuoteOne(t *testing.T) {
	_assert := assert.New(t)

	// right
	var where = map[string]interface{}{
		"code": "sz000001",
		"date": "2020-09-15",
	}

	quote, err := QueryQuoteBaseOne(db.MongoDB, where)
	_assert.Nil(err)
	_assert.Equal("sz000001", quote.Code)
	_assert.Equal("2020-09-15", quote.Date)
	t.Logf("quote: %s", quote.String())

	// no data
	where = map[string]interface{}{
		"code": "no000001",
		"date": "2020-09-15",
	}
	quote, err = QueryQuoteBaseOne(db.MongoDB, where)
	_assert.NotNil(err)
	_assert.Nil(quote)

}

func BenchmarkQueryQuoteBaseOne(b *testing.B) {
	// right
	var where = map[string]interface{}{
		"code": "sz000001",
		"date": "2020-09-15",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QueryQuoteBaseOne(db.MongoDB, where)
	}
}
