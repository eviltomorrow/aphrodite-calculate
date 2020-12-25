package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateBollDay(t *testing.T) {
	_assert := assert.New(t)
	var date = "2020-12-02"
	_, err := CalculateBollDay(date)
	_assert.Nil(err)
}

func TestCalculateBollWeek(t *testing.T) {
	_assert := assert.New(t)
	var date = "2020-12-04"
	_, err := CalculateBollWeek(date)
	_assert.Nil(err)
}
