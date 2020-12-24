package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateMADay(t *testing.T) {
	_assert := assert.New(t)
	var date = "2020-12-02"
	_, err := CalculateMADay(date)
	_assert.Nil(err)
}

func TestCalculateMAWeek(t *testing.T) {
	_assert := assert.New(t)
	var date = "2020-12-04"
	_, err := CalculateMAWeek(date)
	_assert.Nil(err)
}
