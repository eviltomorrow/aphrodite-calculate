package model

import (
	"database/sql"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
)

// SelectTaskRecordMany select task record many
func SelectTaskRecordMany(db db.ExecMySQL, complete bool, offset int64, limit int64) ([]*TaskRecord, error) {
	return nil, nil
}

// DeleteTaskRecordByCodesDate delete task record by codes and date
func DeleteTaskRecordByCodesDate(db db.ExecMySQL, codes []string, date string) (int64, error) {
	return 0, nil
}

// UpdateTaskRecordByID update task record by id
func UpdateTaskRecordByID(db db.ExecMySQL, id int64, record *TaskRecord) (int64, error) {
	return 0, nil
}

// InsertTaskRecordMany insert into task record many
func InsertTaskRecordMany(db db.ExecMySQL, records []*TaskRecord) (int64, error) {
	return 0, nil
}

// TaskRecord record
type TaskRecord struct {
	ID              int64        `json:"id"`   // id
	Name            string       `json:"name"` // 名称
	Code            string       `json:"code"`
	Date            string       `json:"date"`
	Completed       bool         `json:"completed"` // 完成
	Desc            string       `json:"desc"`      // 描述
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}
