package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/eviltomorrow/aphrodite-calculate/db"
)

// SelectTaskRecordManyByCompleted select task record with completed
func SelectTaskRecordManyByCompleted(db db.ExecMySQL, completed bool) ([]*TaskRecord, error) {
	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var _sql = "select id, method, DATE_FORMAT(date,'%Y-%m-%d'), priority, completed, create_timestamp, modify_timestamp from task_record where completed = ? order by date asc,priority asc"
	rows, err := db.QueryContext(ctx, _sql, completed)
	if err != nil {
		return nil, err
	}

	var records = make([]*TaskRecord, 0, 64)
	for rows.Next() {
		var record = TaskRecord{}
		if err := rows.Scan(
			&record.ID,
			&record.Method,
			&record.Date,
			&record.Priority,
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

// SelectTaskRecordManyByDate select task record with date
func SelectTaskRecordManyByDate(db db.ExecMySQL, date string) ([]*TaskRecord, error) {
	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var _sql = "select id, method, DATE_FORMAT(date,'%Y-%m-%d'), priority, completed, create_timestamp, modify_timestamp from task_record where date = ?"
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
			&record.Priority,
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
func UpdateTaskRecordCompleted(db db.ExecMySQL, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), UpdateTimeout)
	defer cannel()

	var fields = make([]string, 0, len(ids))
	var args = make([]interface{}, 0, len(ids))
	for _, id := range ids {
		fields = append(fields, "?")
		args = append(args, id)
	}
	var _sql = fmt.Sprintf("update task_record set completed = true, modify_timestamp = now() where id in (%s)", strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
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
	var args = make([]interface{}, 0, 3*len(records))
	for _, record := range records {
		fields = append(fields, "(?, ?, ?, false, now(), null)")
		args = append(args, record.Method)
		args = append(args, record.Date)
		args = append(args, record.Priority)
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
	TaskRecordFieldPriority        = "priority"
	TaskRecordFieldCompleted       = "completed"
	TaskRecordFieldCreateTimestamp = "create_timestamp"
	TaskRecordFieldModifyTimestamp = "modify_timestamp"
)

var taskRecordFields = []string{
	TaskRecordFieldMethod,
	TaskRecordFieldDate,
	TaskRecordFieldPriority,
	TaskRecordFieldCompleted,
	TaskRecordFieldCreateTimestamp,
	TaskRecordFieldModifyTimestamp,
}

// TaskRecord record
type TaskRecord struct {
	ID              int64        `json:"id"`     // id
	Method          string       `json:"method"` // 方式
	Date            string       `json:"date"`
	Priority        int          `json:"priority"`  // 优先级
	Completed       bool         `json:"completed"` // 完成
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

func (t *TaskRecord) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(t)
	return string(buf)
}
