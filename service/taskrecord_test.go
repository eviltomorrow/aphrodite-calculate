package service

import (
	"testing"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/stretchr/testify/assert"
)

var beginDate = "2020-09-01"
var endDate = "2020-12-10"

func TestBuildTaskRecord(t *testing.T) {
	db.MySQL.Exec("truncate table task_record")

	_assert := assert.New(t)
	begin, err := time.ParseInLocation("2006-01-02", beginDate, time.Local)
	_assert.Nil(err)

	end, err := time.ParseInLocation("2006-01-02", endDate, time.Local)
	_assert.Nil(err)
	err = BuildTaskRecord(begin, end)
	_assert.Nil(err)

	var count int
	for {
		if begin.After(end) {
			break
		}

		switch begin.Weekday() {
		case time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
			count++
		case time.Monday:
			count += 2
		}
		begin = begin.AddDate(0, 0, 1)
	}

	records, err := PollTaskRecord(false)
	_assert.Nil(err)
	_assert.Equal(count, len(records))

	for _, record := range records {
		_assert.Equal(false, record.Completed)
		_assert.Equal(priorityLib[record.Method], record.Priority)
	}

	begin, err = time.ParseInLocation("2006-01-02", beginDate, time.Local)
	_assert.Nil(err)

	end, err = time.ParseInLocation("2006-01-02", endDate, time.Local)
	_assert.Nil(err)
	err = BuildTaskRecord(begin, end)
	_assert.Nil(err)

	records, err = PollTaskRecord(false)
	_assert.Nil(err)
	_assert.Equal(count, len(records))

	for _, record := range records {
		_assert.Equal(false, record.Completed)
		_assert.Equal(priorityLib[record.Method], record.Priority)
	}
}

func TestPollTaskRecord(t *testing.T) {
	db.MySQL.Exec("truncate table task_record")

	_assert := assert.New(t)
	begin, err := time.ParseInLocation("2006-01-02", beginDate, time.Local)
	_assert.Nil(err)

	end, err := time.ParseInLocation("2006-01-02", endDate, time.Local)
	_assert.Nil(err)
	err = BuildTaskRecord(begin, end)
	_assert.Nil(err)

	var count int
	for {
		if begin.After(end) {
			break
		}

		switch begin.Weekday() {
		case time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
			count++
		case time.Monday:
			count += 2
		}
		begin = begin.AddDate(0, 0, 1)
	}

	records, err := PollTaskRecord(false)
	_assert.Nil(err)
	_assert.Equal(count, len(records))

	for _, record := range records {
		_assert.Equal(false, record.Completed)
		_assert.Equal(priorityLib[record.Method], record.Priority)
	}

}

func TestArchiveTaskRecord(t *testing.T) {
	db.MySQL.Exec("truncate table task_record")

	_assert := assert.New(t)
	begin, err := time.ParseInLocation("2006-01-02", beginDate, time.Local)
	_assert.Nil(err)

	end, err := time.ParseInLocation("2006-01-02", endDate, time.Local)
	_assert.Nil(err)
	err = BuildTaskRecord(begin, end)
	_assert.Nil(err)

	var count int
	for {
		if begin.After(end) {
			break
		}

		switch begin.Weekday() {
		case time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
			count++
		case time.Monday:
			count += 2
		}
		begin = begin.AddDate(0, 0, 1)
	}

	records, err := PollTaskRecord(false)
	_assert.Nil(err)
	_assert.Equal(count, len(records))

	for _, record := range records {
		_assert.Equal(false, record.Completed)
		_assert.Equal(priorityLib[record.Method], record.Priority)

		record.NumOfTimes = record.NumOfTimes + 1
		record.Completed = true
		err = ArchiveTaskRecord(record)
		_assert.Nil(err)
	}

	records, err = PollTaskRecord(true)
	_assert.Nil(err)
	_assert.Equal(count, len(records))
	for _, record := range records {
		_assert.Equal(true, record.Completed)
		_assert.Equal(priorityLib[record.Method], record.Priority)
	}
}
