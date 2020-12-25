package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateKDJDay(t *testing.T) {
	_assert := assert.New(t)
	var date = "2020-10-19"
	_, err := CalculateKDJDay(date)
	_assert.Nil(err)
}

func TestCalculateKDJWeek(t *testing.T) {
	_assert := assert.New(t)
	var date = "2020-11-27"
	_, err := CalculateKDJWeek(date)
	_assert.Nil(err)
}

func TestIbuildKDJWeek(t *testing.T) {
	_assert := assert.New(t)
	var date = "2020-10-19"
	kdj, err := buildKDJWeek("sz300999", date)
	_assert.Nil(err)
	t.Logf("%s,\r\n", kdj)
}
