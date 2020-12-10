package model

import (
	"testing"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

func TestQueryQuoteBaseCurrentCodeLimit2(t *testing.T) {
	_assert := assert.New(t)

	var code = "sz000001"
	var date = "2020-12-04"
	quotes, err := SelectQuoteBaseCurrentCodeLimit2(db.MongoDB, code, date)
	_assert.Nil(err)
	_assert.Equal(2, len(quotes))

	for _, quote := range quotes {
		t.Logf("Quote: %s\r\n", quote.String())
	}

	code = "sz000001"
	date = "2020-12-08"
	quotes, err = SelectQuoteBaseCurrentCodeLimit2(db.MongoDB, code, date)
	_assert.Nil(err)
	_assert.Equal(1, len(quotes))

	code = "sz000001"
	date = "2020-12-09"
	quotes, err = SelectQuoteBaseCurrentCodeLimit2(db.MongoDB, code, date)
	_assert.Nil(err)
	_assert.Equal(0, len(quotes))

	code = "sz000001"
	date = "2020-12-31"
	quotes, err = SelectQuoteBaseCurrentCodeLimit2(db.MongoDB, code, date)
	_assert.Nil(err)
	_assert.Equal(0, len(quotes))

}

func BenchmarkQueryQuoteBaseOne(b *testing.B) {
	// right
	var code = "sz000001"
	var date = "2020-12-05"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SelectQuoteBaseCurrentCodeLimit2(db.MongoDB, code, date)
	}
}
