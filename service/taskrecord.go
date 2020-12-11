package service

import (
	"fmt"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/eviltomorrow/aphrodite-calculate/model"
)

var standardTaskMethodLib = []string{
	"SYNC_QUOTEDAY",
	"SYNC_QUOTEWEEK",
}

// BuildTaskRecord build task record
func BuildTaskRecord(begin, end time.Time) error {
	if begin.Format("2006-01-02") != end.Format("2006-01-02") && begin.After(end) {
		return fmt.Errorf("Invalid date, begin: %v, end: %v", begin, end)
	}

	var cache = make([]*model.TaskRecord, 0, 64)
	for {
		if begin.After(end) {
			break
		}

		var current = begin.Format("2006-01-02")
		begin = begin.AddDate(0, 0, 1)

		records, err := model.SelectTaskRecordMany(db.MySQL, current)
		if err != nil {
			return err
		}

	loop:
		for _, method := range standardTaskMethodLib {
			for _, record := range records {
				if method == record.Method {
					continue loop
				}
			}

			var record = &model.TaskRecord{
				Method:    method,
				Date:      current,
				Completed: false,
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
	}

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
	return nil
}

// PollTaskRecord poll task record
func PollTaskRecord(date string) ([]*model.TaskRecord, error) {
	return model.SelectTaskRecordMany(db.MySQL, date)
}

// ArchiveTaskRecord archive task record
func ArchiveTaskRecord(ids []int64) error {
	tx, err := db.MySQL.Begin()
	if err != nil {
		return err
	}

	_, err = model.UpdateTaskRecordCompleted(tx, ids)
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
