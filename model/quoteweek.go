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

// SelectQuoteWeekByCodeDate select quoteday
func SelectQuoteWeekByCodeDate(db db.ExecMySQL, codes []string, date string) ([]QuoteWeek, error) {
	if len(codes) == 0 {
		return []QuoteWeek{}, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var fields = make([]string, 0, len(codes))
	var args = make([]interface{}, 0, len(codes))
	for _, code := range codes {
		fields = append(fields, "?")
		args = append(args, code)
	}
	args = append(args, date)
	var _sql = fmt.Sprintf("select id, code, open, close, high, low, volume, account, date_begin, date_end, week_of_year, create_timestamp, modify_timestamp from quote_week where code in (%s) and date_end = ?", strings.Join(fields, ","))

	rows, err := db.QueryContext(ctx, _sql, args...)
	if err != nil {
		return nil, err
	}

	var quotes = make([]QuoteWeek, 0, len(codes))
	for rows.Next() {
		var quote = QuoteWeek{}
		if err := rows.Scan(
			&quote.ID,
			&quote.Code,
			&quote.Open,
			&quote.Close,
			&quote.High,
			&quote.Low,
			&quote.Volume,
			&quote.Account,
			&quote.DateBegin,
			&quote.DateEnd,
			&quote.WeekOfYear,
			&quote.CreateTimestamp,
			&quote.ModifyTimestamp,
		); err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return quotes, nil
}

// DeleteQuoteWeekByCodesDate delete quoteweek
func DeleteQuoteWeekByCodesDate(db db.ExecMySQL, codes []string, date string) (int64, error) {
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

	var _sql = fmt.Sprintf("delete from quote_week where code in (%s) and date_end = ?", strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertQuoteWeekMany batch insert quoteweek for mysql
func InsertQuoteWeekMany(db db.ExecMySQL, quotes []*QuoteWeek) (int64, error) {
	if len(quotes) == 0 {
		return 0, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), InsertTimeout)
	defer cannel()

	var fields = make([]string, 0, len(quotes))
	var args = make([]interface{}, 0, 3*len(quotes))
	for _, quote := range quotes {
		fields = append(fields, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now(), null)")
		args = append(args, quote.Code)
		args = append(args, quote.Open)
		args = append(args, quote.Close)
		args = append(args, quote.High)
		args = append(args, quote.Low)
		args = append(args, quote.Volume)
		args = append(args, quote.Account)
		args = append(args, quote.DateBegin)
		args = append(args, quote.DateEnd)
		args = append(args, quote.WeekOfYear)
	}

	var _sql = fmt.Sprintf("insert into quote_week (%s) values %s", strings.Join(quoteWeekFields, ","), strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//
const (
	QuoteWeekFieldCode            = "code"
	QuoteWeekFieldOpen            = "open"
	QuoteWeekFieldClose           = "close"
	QuoteWeekFieldHigh            = "high"
	QuoteWeekFieldLow             = "low"
	QuoteWeekFieldVolume          = "volume"
	QuoteWeekFieldAccount         = "account"
	QuoteWeekFieldDateBegin       = "date_begin"
	QuoteWeekFieldDateEnd         = "date_end"
	QuoteWeekFieldWeekOfYear      = "week_of_year"
	QuoteWeekFieldCreateTimestamp = "create_timestamp"
	QuoteWeekFieldModifyTimestamp = "modify_timestamp"
)

var quoteWeekFields = []string{
	QuoteWeekFieldCode,
	QuoteWeekFieldOpen,
	QuoteWeekFieldClose,
	QuoteWeekFieldHigh,
	QuoteWeekFieldLow,
	QuoteWeekFieldVolume,
	QuoteWeekFieldAccount,
	QuoteWeekFieldDateBegin,
	QuoteWeekFieldDateEnd,
	QuoteWeekFieldWeekOfYear,
	QuoteWeekFieldCreateTimestamp,
	QuoteWeekFieldModifyTimestamp,
}

// QuoteWeek quote week
type QuoteWeek struct {
	ID              int64        `json:"id"`
	Code            string       `json:"code"`
	Open            float64      `json:"open"`
	Close           float64      `json:"close"`
	High            float64      `json:"high"`
	Low             float64      `json:"low"`
	Volume          int64        `json:"volume"`
	Account         float64      `json:"account"`
	DateBegin       time.Time    `json:"date_begin"`
	DateEnd         time.Time    `json:"date_end"`
	WeekOfYear      int          `json:"week_of_year"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

func (q *QuoteWeek) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(q)
	return string(buf)
}
