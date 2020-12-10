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

// SelectQuoteDayByCodeDate select quoteday
func SelectQuoteDayByCodeDate(db db.ExecMySQL, code string, begin, end string) ([]*QuoteDay, error) {
	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var _sql = fmt.Sprintf("select id, code, open, close, high, low, volume, account, date, day_of_year, create_timestamp, modify_timestamp from quote_day where code = ? and date between ? and ?")

	rows, err := db.QueryContext(ctx, _sql, code, begin, end)
	if err != nil {
		return nil, err
	}

	var quotes = make([]*QuoteDay, 0, 16)
	for rows.Next() {
		var quote = QuoteDay{}
		if err := rows.Scan(
			&quote.ID,
			&quote.Code,
			&quote.Open,
			&quote.Close,
			&quote.High,
			&quote.Low,
			&quote.Volume,
			&quote.Account,
			&quote.Date,
			&quote.DayOfYear,
			&quote.CreateTimestamp,
			&quote.ModifyTimestamp,
		); err != nil {
			return nil, err
		}
		quotes = append(quotes, &quote)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return quotes, nil
}

// DeleteQuoteDayByCodesDate delete quoteday
func DeleteQuoteDayByCodesDate(db db.ExecMySQL, codes []string, date string) (int64, error) {
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

	var _sql = fmt.Sprintf("delete from quote_day where code in (%s) and date = ?", strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertQuoteDayMany batch insert quoteday for mysql
func InsertQuoteDayMany(db db.ExecMySQL, quotes []*QuoteDay) (int64, error) {
	if len(quotes) == 0 {
		return 0, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), InsertTimeout)
	defer cannel()

	var fields = make([]string, 0, len(quotes))
	var args = make([]interface{}, 0, 9*len(quotes))
	for _, quote := range quotes {
		fields = append(fields, "(?, ?, ?, ?, ?, ?, ?, ?, ?, now(), null)")
		args = append(args, quote.Code)
		args = append(args, quote.Open)
		args = append(args, quote.Close)
		args = append(args, quote.High)
		args = append(args, quote.Low)
		args = append(args, quote.Volume)
		args = append(args, quote.Account)
		args = append(args, quote.Date)
		args = append(args, quote.DayOfYear)
	}

	var _sql = fmt.Sprintf("insert into quote_day (%s) values %s", strings.Join(quoteDayFields, ","), strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//
const (
	QuoteDayFieldCode            = "code"
	QuoteDayFieldOpen            = "open"
	QuoteDayFieldClose           = "close"
	QuoteDayFieldHigh            = "high"
	QuoteDayFieldLow             = "low"
	QuoteDayFieldVolume          = "volume"
	QuoteDayFieldAccount         = "account"
	QuoteDayFieldDate            = "date"
	QuoteDayFieldDayOfYear       = "day_of_year"
	QuoteDayFieldCreateTimestamp = "create_timestamp"
	QuoteDayFieldModifyTimestamp = "modify_timestamp"
)

var quoteDayFields = []string{
	QuoteDayFieldCode,
	QuoteDayFieldOpen,
	QuoteDayFieldClose,
	QuoteDayFieldHigh,
	QuoteDayFieldLow,
	QuoteDayFieldVolume,
	QuoteDayFieldAccount,
	QuoteDayFieldDate,
	QuoteDayFieldDayOfYear,
	QuoteDayFieldCreateTimestamp,
	QuoteDayFieldModifyTimestamp,
}

// QuoteDay quote day
type QuoteDay struct {
	ID              int64        `json:"id"`
	Code            string       `json:"code"`
	Open            float64      `json:"open"`
	Close           float64      `json:"close"`
	High            float64      `json:"high"`
	Low             float64      `json:"low"`
	Volume          int64        `json:"volume"`
	Account         float64      `json:"account"`
	Date            time.Time    `json:"date"`
	DayOfYear       int          `json:"day_of_year"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

func (q *QuoteDay) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(q)
	return string(buf)
}
