package service

import (
	"fmt"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/eviltomorrow/aphrodite-calculate/model"
)

//
const (
	SyncQuoteDay  = "SYNC_QUOTEDAY"
	SyncQuoteWeek = "SYNC_QUOTEWEEK"
)

var priorityLib = map[string]int{
	SyncQuoteDay:  0,
	SyncQuoteWeek: 1,
}

// BuildTaskRecord build task record
func BuildTaskRecord(begin, end time.Time) error {
	if begin.Format("2006-01-02") != end.Format("2006-01-02") && begin.After(end) {
		return fmt.Errorf("Invalid date, begin: %v, end: %v", begin, end)
	}

	var cache = make([]*model.TaskRecord, 0, 64)
	var methods = make([]string, 0, 8)
	for {
		if begin.After(end) {
			break
		}
		var current = begin.Format("2006-01-02")

		records, err := model.SelectTaskRecordManyByDate(db.MySQL, current)
		if err != nil {
			return err
		}

		switch begin.Weekday() {
		case time.Monday, time.Tuesday, time.Wednesday, time.Thursday:
			methods = append(methods, SyncQuoteDay)
		case time.Friday:
			methods = append(methods, SyncQuoteDay)
			methods = append(methods, SyncQuoteWeek)
		default:
		}

	loop:
		for _, method := range methods {
			for _, record := range records {
				if record.Method == method {
					continue loop
				}
			}

			var record = &model.TaskRecord{
				Method:    method,
				Date:      current,
				Completed: false,
				Priority:  priorityLib[method],
			}
			cache = append(cache, record)
		}

		if len(cache) > 60 {
			tx, err := db.MySQL.Begin()
			if err != nil {
				return err
			}
			if _, err = model.InsertTaskRecordMany(tx, cache); err != nil {
				tx.Rollback()
				return err
			}
			if err = tx.Commit(); err != nil {
				tx.Rollback()
				return err
			}
			cache = cache[:0]
		}

		methods = methods[:0]
		begin = begin.AddDate(0, 0, 1)
	}

	if len(cache) > 0 {
		tx, err := db.MySQL.Begin()
		if err != nil {
			return err
		}
		if _, err = model.InsertTaskRecordMany(tx, cache); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

// ArchiveTaskRecord archive task record
func ArchiveTaskRecord(record *model.TaskRecord) error {
	if record == nil {
		return nil
	}

	tx, err := db.MySQL.Begin()
	if err != nil {
		return err
	}

	_, err = model.UpdateTaskRecord(tx, record, record.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// PollTaskRecord poll uncompleted task record
func PollTaskRecord(completed bool) ([]*model.TaskRecord, error) {
	records, err := model.SelectTaskRecordManyByCompleted(db.MySQL, completed)
	if err != nil {
		return nil, err
	}

	var cache = make([]*model.TaskRecord, 0, len(records))
	for _, record := range records {
		if record.NumOfTimes < 14 {
			cache = append(cache, record)
		}
	}

	return cache, nil
}
