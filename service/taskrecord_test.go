package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildTaskRecord(t *testing.T) {
	_assert := assert.New(t)
	begin, err := time.ParseInLocation("2006-01-02", "2020-12-01", time.Local)
	_assert.Nil(err)

	end, err := time.ParseInLocation("2006-01-02", "2020-12-10", time.Local)
	_assert.Nil(err)
	err = BuildTaskRecord(begin, end)
	_assert.Nil(err)
}

func BenchmarkBuildTaskRecord(b *testing.B) {
	begin, err := time.ParseInLocation("2006-01-02", "2020-12-01", time.Local)
	if err != nil {
		b.Fatalf("Error: %v\r\n", err)
	}

	end, err := time.ParseInLocation("2006-01-02", "2020-12-10", time.Local)
	if err != nil {
		b.Fatalf("Error: %v\r\n", err)
	}

	for i := 0; i < b.N; i++ {
		BuildTaskRecord(begin, end)
	}

}
