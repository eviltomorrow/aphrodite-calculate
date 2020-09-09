package model

import (
	"testing"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

func TestQueryQuoteOne(t *testing.T) {
	_assert := assert.New(t)

	var where = map[string]interface{}{
		"code": "sz000055",
		"date": "2020-06-03",
	}

	quote, err := QueryQuoteOne(db.MongoDB, where)
	_assert.Nil(err)
	t.Log(quote.String())
}
