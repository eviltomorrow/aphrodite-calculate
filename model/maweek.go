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

// DeleteMAWeekByCodesDate delete ma day by code
func DeleteMAWeekByCodesDate(db db.ExecMySQL, codes []string, date string) (int64, error) {
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

	var _sql = fmt.Sprintf("delete from ma_week where code in (%s) and date = ?", strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertMAWeekMany batch insert ma day
func InsertMAWeekMany(db db.ExecMySQL, mas []*MAWeek) (int64, error) {
	if len(mas) == 0 {
		return 0, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), InsertTimeout)
	defer cannel()

	var fields = make([]string, 0, len(mas))
	var args = make([]interface{}, 0, 10*len(mas))
	for _, ma := range mas {
		fields = append(fields, "(?, ?, ?, ?, ?, ?, ?, now(), null)")
		args = append(args, ma.Code)
		args = append(args, ma.M5)
		args = append(args, ma.M10)
		args = append(args, ma.M20)
		args = append(args, ma.M30)
		args = append(args, ma.Date)
		args = append(args, ma.WeekOfYear)
	}

	var _sql = fmt.Sprintf("insert into ma_week (%s) values %s", strings.Join(maWeekFields, ","), strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//
const (
	MAWeekFieldCode            = "code"
	MAWeekFieldM5              = "m5"
	MAWeekFieldM10             = "m10"
	MAWeekFieldM20             = "m20"
	MAWeekFieldM30             = "m30"
	MAWeekFieldDate            = "date"
	MAWeekFieldWeekOfYear      = "week_of_year"
	MAWeekFieldCreateTimestamp = "create_timestamp"
	MAWeekFieldModifyTimestamp = "modify_timestamp"
)

var maWeekFields = []string{
	MAWeekFieldCode,
	MAWeekFieldM5,
	MAWeekFieldM10,
	MAWeekFieldM20,
	MAWeekFieldM30,
	MAWeekFieldDate,
	MAWeekFieldWeekOfYear,
	MAWeekFieldCreateTimestamp,
	MAWeekFieldModifyTimestamp,
}

// MAWeek ma week
type MAWeek struct {
	ID              int64        `json:"id"`
	Code            string       `json:"code"`
	M5              float64      `json:"m5"`
	M10             float64      `json:"m10"`
	M20             float64      `json:"m20"`
	M30             float64      `json:"m30"`
	Date            time.Time    `json:"date"`
	WeekOfYear      int          `json:"week_of_year"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

func (m *MAWeek) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(m)
	return string(buf)
}
