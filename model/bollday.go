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

// DeleteBollDayByCodesDate delete Boll day by code
func DeleteBollDayByCodesDate(db db.ExecMySQL, codes []string, date string) (int64, error) {
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

	var _sql = fmt.Sprintf("delete from boll_day where code in (%s) and date = ?", strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertBollDayMany batch insert Boll day
func InsertBollDayMany(db db.ExecMySQL, bolls []*BollDay) (int64, error) {
	if len(bolls) == 0 {
		return 0, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), InsertTimeout)
	defer cannel()

	var fields = make([]string, 0, len(bolls))
	var args = make([]interface{}, 0, 10*len(bolls))
	for _, boll := range bolls {
		fields = append(fields, "(?, ?, ?, ?, ?, ?, now(), null)")
		args = append(args, boll.Code)
		args = append(args, boll.UP)
		args = append(args, boll.MB)
		args = append(args, boll.DN)
		args = append(args, boll.Date)
		args = append(args, boll.DayOfYear)
	}

	var _sql = fmt.Sprintf("insert into boll_day (%s) values %s", strings.Join(bollDayFields, ","), strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//
const (
	BollDayFieldCode            = "code"
	BollDayFieldUP              = "up"
	BollDayFieldMB              = "mb"
	BollDayFieldDN              = "dn"
	BollDayFieldDate            = "date"
	BollDayFieldDayOfYear       = "day_of_year"
	BollDayFieldCreateTimestamp = "create_timestamp"
	BollDayFieldModifyTimestamp = "modify_timestamp"
)

var bollDayFields = []string{
	BollDayFieldCode,
	BollDayFieldUP,
	BollDayFieldMB,
	BollDayFieldDN,
	BollDayFieldDate,
	BollDayFieldDayOfYear,
	BollDayFieldCreateTimestamp,
	BollDayFieldModifyTimestamp,
}

// BollDay Boll day
type BollDay struct {
	ID              int64        `json:"id"`
	Code            string       `json:"code"`
	UP              float64      `json:"up"`
	MB              float64      `json:"mb"`
	DN              float64      `json:"dn"`
	Date            time.Time    `json:"date"`
	DayOfYear       int          `json:"day_of_year"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

func (b *BollDay) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(b)
	return string(buf)
}
