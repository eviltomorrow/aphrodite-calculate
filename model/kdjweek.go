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

// DeleteKDJWeekByCodesDate delete KDJ day by code
func DeleteKDJWeekByCodesDate(db db.ExecMySQL, codes []string, date string) (int64, error) {
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

	var _sql = fmt.Sprintf("delete from kdj_week where code in (%s) and date = ?", strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertKDJWeekMany batch insert KDJ day
func InsertKDJWeekMany(db db.ExecMySQL, kdjs []*KDJWeek) (int64, error) {
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
		args = append(args, kdj.WeekOfYear)
	}

	var _sql = fmt.Sprintf("insert into kdj_week (%s) values %s", strings.Join(kdjWeekFields, ","), strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//
const (
	KDJWeekFieldCode            = "code"
	KDJWeekFieldUP              = "up"
	KDJWeekFieldMB              = "mb"
	KDJWeekFieldDN              = "dn"
	KDJWeekFieldDate            = "date"
	KDJWeekFieldWeekOfYear      = "week_of_year"
	KDJWeekFieldCreateTimestamp = "create_timestamp"
	KDJWeekFieldModifyTimestamp = "modify_timestamp"
)

var kdjWeekFields = []string{
	KDJWeekFieldCode,
	KDJWeekFieldUP,
	KDJWeekFieldMB,
	KDJWeekFieldDN,
	KDJWeekFieldDate,
	KDJWeekFieldWeekOfYear,
	KDJWeekFieldCreateTimestamp,
	KDJWeekFieldModifyTimestamp,
}

// KDJWeek KDJ week
type KDJWeek struct {
	ID              int64        `json:"id"`
	Code            string       `json:"code"`
	K               float64      `json:"k"`
	D               float64      `json:"d"`
	J               float64      `json:"j"`
	Date            time.Time    `json:"date"`
	WeekOfYear      int          `json:"week_of_year"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

func (k *KDJWeek) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(k)
	return string(buf)
}
