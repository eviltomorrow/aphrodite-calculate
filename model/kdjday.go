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

// SelectKDJDayByCodeDateLimit select kdj day
func SelectKDJDayByCodeDateLimit(db db.ExecMySQL, code string, date string, limit int64) ([]*KDJDay, error) {
	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var _sql = fmt.Sprintf("select k, d, j from kdj_day where code = ? and date < ? order by date desc limit ?")

	rows, err := db.QueryContext(ctx, _sql, code, date, limit)
	if err != nil {
		return nil, err
	}

	var kdjs = make([]*KDJDay, 0, limit)
	for rows.Next() {
		var kdj = KDJDay{}
		if err := rows.Scan(
			&kdj.K,
			&kdj.D,
			&kdj.J,
		); err != nil {
			return nil, err
		}
		kdjs = append(kdjs, &kdj)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return kdjs, nil
}

// DeleteKDJDayByCodesDate delete KDJ day by code
func DeleteKDJDayByCodesDate(db db.ExecMySQL, codes []string, date string) (int64, error) {
	if len(codes) == 0 {
		return 0, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), DeleteTimeout)
	defer cannel()

	var fields = make([]string, 0, len(codes))
	var args = make([]interface{}, 0, len(codes)+1)
	for _, code := range codes {
		fields = append(fields, "?")
		args = append(args, code)
	}
	args = append(args, date)

	var _sql = fmt.Sprintf("delete from kdj_day where code in (%s) and date = ?", strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertKDJDayMany batch insert KDJ day
func InsertKDJDayMany(db db.ExecMySQL, kdjs []*KDJDay) (int64, error) {
	if len(kdjs) == 0 {
		return 0, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), InsertTimeout)
	defer cannel()

	var fields = make([]string, 0, len(kdjs))
	var args = make([]interface{}, 0, 10*len(kdjs))
	for _, kdj := range kdjs {
		fields = append(fields, "(?, ?, ?, ?, ?, ?, now(), null)")
		args = append(args, kdj.Code)
		args = append(args, kdj.K)
		args = append(args, kdj.D)
		args = append(args, kdj.J)
		args = append(args, kdj.Date)
		args = append(args, kdj.DayOfYear)
	}

	var _sql = fmt.Sprintf("insert into kdj_day (%s) values %s", strings.Join(kdjDayFields, ","), strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//
const (
	KDJDayFieldCode            = "code"
	KDJDayFieldK               = "k"
	KDJDayFieldD               = "d"
	KDJDayFieldJ               = "j"
	KDJDayFieldDate            = "date"
	KDJDayFieldDayOfYear       = "day_of_year"
	KDJDayFieldCreateTimestamp = "create_timestamp"
	KDJDayFieldModifyTimestamp = "modify_timestamp"
)

var kdjDayFields = []string{
	KDJDayFieldCode,
	KDJDayFieldK,
	KDJDayFieldD,
	KDJDayFieldJ,
	KDJDayFieldDate,
	KDJDayFieldDayOfYear,
	KDJDayFieldCreateTimestamp,
	KDJDayFieldModifyTimestamp,
}

// KDJDay KDJ day
type KDJDay struct {
	ID              int64        `json:"id"`
	Code            string       `json:"code"`
	K               float64      `json:"k"`
	D               float64      `json:"d"`
	J               float64      `json:"j"`
	Date            time.Time    `json:"date"`
	DayOfYear       int          `json:"day_of_year"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

func (k *KDJDay) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(k)
	return string(buf)
}
