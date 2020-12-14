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

		records, err := model.SelectTaskRecordManyByDate(db.MySQL, begin.Format("2006-01-02"))
		if err != nil {
			return err
		}

		switch begin.Weekday() {
		case time.Monday, time.Tuesday, time.Wednesday, time.Thursday:
			if len(records) == 1 && records[0].Method == SyncQuoteDay {
				break
			}
			if len(records) == 0 {
				cache = append(cache, &model.TaskRecord{
					Method:    SyncQuoteDay,
					Date:      begin.Format("2006-01-02"),
					Completed: false,
				})
			}
		case time.Friday:

		default:
		}
	}
	// 	switch begin.Weekday() {
	// 	case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:

	// 	// case :
	// 	default:
	// 	}
	// }

	// 	var current = begin.Format("2006-01-02")
	// 	begin = begin.AddDate(0, 0, 1)

	// 	records, err := model.SelectTaskRecordManyByDate(db.MySQL, current)
	// 	if err != nil {
	// 		return err
	// 	}

	// loop:
	// 	for _, method := range standardTaskMethodLib {
	// 		for _, record := range records {
	// 			if method == record.Method {
	// 				continue loop
	// 			}
	// 		}

	// 		var record = &model.TaskRecord{
	// 			Method:    method,
	// 			Date:      current,
	// 			Completed: false,
	// 		}
	// 		cache = append(cache, record)
	// 	}

	// 	if len(cache) > 60 {
	// 		tx, err := db.MySQL.Begin()
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if _, err = model.InsertTaskRecordMany(tx, cache); err != nil {
	// 			tx.Rollback()
	// 			return err
	// 		}
	// 		if err = tx.Commit(); err != nil {
	// 			tx.Rollback()
	// 			return err
	// 		}
	// 		cache = cache[:0]
	// 	}
	// }

	// tx, err := db.MySQL.Begin()
	// if err != nil {
	// 	return err
	// }
	// if _, err = model.InsertTaskRecordMany(tx, cache); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// if err = tx.Commit(); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	return nil
}

// PollUncompletedTaskRecord poll uncompleted task record
func PollUncompletedTaskRecord(completed bool) ([]*model.TaskRecord, error) {
	return model.SelectTaskRecordManyByCompleted(db.MySQL, completed)
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
