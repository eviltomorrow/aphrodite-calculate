package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateKDJDay(t *testing.T) {
	_assert := assert.New(t)
	var date = "2020-10-15"
	_, err := CalculateKDJDay(date)
	_assert.Nil(err)
}

func TestIbuildKDJDay(t *testing.T) {
	_assert := assert.New(t)
	kdj, err := buildKDJDay("sz300999", "2020-10-15")
	_assert.Nil(err)
	t.Logf("%s\r\n", kdj)
}
