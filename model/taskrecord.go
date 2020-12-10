package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	jsoniter "github.com/json-iterator/go"
)

// SelectTaskRecordMany select task record many
func SelectTaskRecordMany(db db.ExecMySQL, date string) ([]*TaskRecord, error) {
	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var _sql = "select id, method, date, completed, create_timestamp, modify_timestamp from task_record where date = ?"
	rows, err := db.QueryContext(ctx, _sql, date)
	if err != nil {
		return nil, err
	}

	var records = make([]*TaskRecord, 0, 16)
	for rows.Next() {
		var record = TaskRecord{}
		if err := rows.Scan(
			&record.ID,
			&record.Method,
			&record.Date,
			&record.Completed,
			&record.CreateTimestamp,
			&record.ModifyTimestamp,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

// UpdateTaskRecordCompleted update task record by id
func UpdateTaskRecordCompleted(db db.ExecMySQL, id int64, completed bool) (int64, error) {
	ctx, cannel := context.WithTimeout(context.Background(), UpdateTimeout)
	defer cannel()

	var _sql = "update task_record set completed = ?, modify_timestamp = now() where id = ?"
	result, err := db.ExecContext(ctx, _sql, completed, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertTaskRecordMany insert into task record many
func InsertTaskRecordMany(db db.ExecMySQL, records []*TaskRecord) (int64, error) {
	if len(records) == 0 {
		return 0, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), InsertTimeout)
	defer cannel()

	var fields = make([]string, 0, len(records))
	var args = make([]interface{}, 0, 2*len(records))
	for _, record := range records {
		fields = append(fields, "(?, ?, false, now(), null)")
		args = append(args, record.Method)
		args = append(args, record.Date)
	}

	var _sql = fmt.Sprintf("insert into task_record (%s) values %s", strings.Join(taskRecordFields, ","), strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//
const (
	TaskRecordFieldID              = "id"
	TaskRecordFieldMethod          = "method"
	TaskRecordFieldDate            = "date"
	TaskRecordFieldCompleted       = "completed"
	TaskRecordFieldCreateTimestamp = "create_timestamp"
	TaskRecordFieldModifyTimestamp = "modify_timestamp"
)

var taskRecordFields = []string{
	TaskRecordFieldMethod,
	TaskRecordFieldDate,
	TaskRecordFieldCompleted,
	TaskRecordFieldCreateTimestamp,
	TaskRecordFieldModifyTimestamp,
}

// TaskRecord record
type TaskRecord struct {
	ID              int64        `json:"id"`     // id
	Method          string       `json:"method"` // 方式
	Date            string       `json:"date"`
	Completed       bool         `json:"completed"` // 完成
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

func (t *TaskRecord) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(t)
	return string(buf)
}
