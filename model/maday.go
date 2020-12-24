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

// DeleteMADayByCodesDate delete ma day by code
func DeleteMADayByCodesDate(db db.ExecMySQL, codes []string, date string) (int64, error) {
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

	var _sql = fmt.Sprintf("delete from ma_day where code in (%s) and date = ?", strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertMADayMany batch insert ma day
func InsertMADayMany(db db.ExecMySQL, mas []*MADay) (int64, error) {
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
		args = append(args, ma.DayOfYear)
	}

	var _sql = fmt.Sprintf("insert into ma_day (%s) values %s", strings.Join(maDayFields, ","), strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//
const (
	MADayFieldCode            = "code"
	MADayFieldM5              = "m5"
	MADayFieldM10             = "m10"
	MADayFieldM20             = "m20"
	MADayFieldM30             = "m30"
	MADayFieldDate            = "date"
	MADayFieldDayOfYear       = "day_of_year"
	MADayFieldCreateTimestamp = "create_timestamp"
	MADayFieldModifyTimestamp = "modify_timestamp"
)

var maDayFields = []string{
	MADayFieldCode,
	MADayFieldM5,
	MADayFieldM10,
	MADayFieldM20,
	MADayFieldM30,
	MADayFieldDate,
	MADayFieldDayOfYear,
	MADayFieldCreateTimestamp,
	MADayFieldModifyTimestamp,
}

// MADay ma day
type MADay struct {
	ID              int64        `json:"id"`
	Code            string       `json:"code"`
	M5              float64      `json:"m5"`
	M10             float64      `json:"m10"`
	M20             float64      `json:"m20"`
	M30             float64      `json:"m30"`
	Date            time.Time    `json:"date"`
	DayOfYear       int          `json:"day_of_year"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

func (m *MADay) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(m)
	return string(buf)
}
