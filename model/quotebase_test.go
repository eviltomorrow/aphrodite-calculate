package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eviltomorrow/aphrodite-calculate/db"
)

func TestQueryQuoteBaseCurrentCodeLimit2(t *testing.T) {
	_assert := assert.New(t)

	var code = "sz000001"
	var date = "2020-09-05"
	quotes, err := SelectQuoteBaseCurrentCodeLimit2(db.MongoDB, code, date)
	_assert.Nil(err)
	_assert.Equal(2, len(quotes))

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

	code = "sz000001"
	date = "2020"
	quotes, err = SelectQuoteBaseCurrentCodeLimit2(db.MongoDB, code, date)
	_assert.NotNil(err)
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
