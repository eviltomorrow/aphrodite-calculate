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

func TestPollUncompletedTaskRecord(t *testing.T) {
	_assert := assert.New(t)
	records, err := PollUncompletedTaskRecord(false)
	_assert.Nil(err)
	_assert.Equal(2, len(records))

	records, err = PollUncompletedTaskRecord(false)
	_assert.Nil(err)
	_assert.Equal(0, len(records))
}

func TestArchiveTaskRecord(t *testing.T) {
	_assert := assert.New(t)
	records, err := PollUncompletedTaskRecord(false)
	_assert.Nil(err)

	var ids = make([]int64, 0, len(records))
	for _, record := range records {
		ids = append(ids, record.ID)
	}

	err = ArchiveTaskRecord(ids)
	_assert.Nil(err)
}
