package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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
	if len(records) == 0 {
		return 0, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), InsertTimeout)
	defer cannel()

	var fields = make([]string, 0, len(records))
	var args = make([]interface{}, 0, 3*len(records))
	for _, record := range records {
		fields = append(fields, "(?, ?, ?, false, '', now(), null)")
		args = append(args, record.Name)
		args = append(args, record.Code)
		args = append(args, record.Date)
	}

	var _sql = fmt.Sprintf("insert into task_record (%s) values %s", strings.Join(taskRecordFields, ","), strings.Join(fields, ","))
	fmt.Println(_sql)
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//
const (
	TaskRecordFieldID              = "id"
	TaskRecordFieldName            = "name"
	TaskRecordFieldCode            = "code"
	TaskRecordFieldDate            = "date"
	TaskRecordFieldCompleted       = "completed"
	TaskRecordFieldMsg             = "msg"
	TaskRecordFieldCreateTimestamp = "create_timestamp"
	TaskRecordFieldModifyTimestamp = "modify_timestamp"
)

var taskRecordFields = []string{
	TaskRecordFieldName,
	TaskRecordFieldCode,
	TaskRecordFieldDate,
	TaskRecordFieldCompleted,
	TaskRecordFieldMsg,
	TaskRecordFieldCreateTimestamp,
	TaskRecordFieldModifyTimestamp,
}

// TaskRecord record
type TaskRecord struct {
	ID              int64        `json:"id"`   // id
	Name            string       `json:"name"` // 名称
	Code            string       `json:"code"`
	Date            string       `json:"date"`
	Completed       bool         `json:"completed"` // 完成
	Msg             string       `json:"msg"`       // 描述
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}
