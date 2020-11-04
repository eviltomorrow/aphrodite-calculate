package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// SelectQuoteWeekLatestByCodeDate select quoteweek
func SelectQuoteWeekLatestByCodeDate(db *sql.DB, code string, date string, size int) ([]QuoteWeek, error) {
	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var _sql = "select id, code, open, close, high, low, volume, account, date_begin, date_end, week_of_year, create_timestamp, modify_timestamp from quote_week where code =? and date = ? order by date desc limit ?"
	rows, err := db.QueryContext(ctx, _sql, code, date, size)
	if err != nil {
		return nil, err
	}

	var quotes = make([]QuoteWeek, 0, size)
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

// SelectQuoteWeekByCodeDate select quoteday
func SelectQuoteWeekByCodeDate(db *sql.DB, codes []string, date string) ([]QuoteWeek, error) {
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
	var _sql = fmt.Sprintf("select id, code, open, close, high, low, volume, account, date_begin, date_end, week_of_year, create_timestamp, modify_timestamp from quote_week where code in (%s) and date = ?", strings.Join(fields, ","))

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

// DeleteQuoteWeekByCodeDate delete quoteweek
func DeleteQuoteWeekByCodeDate(db *sql.DB, code string, date string) (int64, error) {
	ctx, cannel := context.WithTimeout(context.Background(), DeleteTimeout)
	defer cannel()

	var _sql = "delete from quote_day where code = ? and date = ?"
	result, err := db.ExecContext(ctx, _sql, code, date)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertQuoteWeekMany batch insert quoteweek for mysql
func InsertQuoteWeekMany(db *sql.DB, quotes []*QuoteWeek) (int64, error) {
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

	var _sql = fmt.Sprintf("insert into quote_week (%s) values %s", strings.Join(quoteWeekFeilds, ","), strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//
const (
	QuoteWeekFeildCode            = "code"
	QuoteWeekFeildOpen            = "open"
	QuoteWeekFeildClose           = "close"
	QuoteWeekFeildHigh            = "high"
	QuoteWeekFeildLow             = "low"
	QuoteWeekFeildVolume          = "volume"
	QuoteWeekFeildAccount         = "account"
	QuoteWeekFeildDateBegin       = "date_begin"
	QuoteWeekFeildDateEnd         = "date_end"
	QuoteWeekFeildWeekOfYear      = "week_of_year"
	QuoteWeekFeildCreateTimestamp = "create_timestamp"
	QuoteWeekFeildModifyTimestamp = "modify_timestamp"
)

var quoteWeekFeilds = []string{
	QuoteWeekFeildCode,
	QuoteWeekFeildOpen,
	QuoteWeekFeildClose,
	QuoteWeekFeildHigh,
	QuoteWeekFeildLow,
	QuoteWeekFeildVolume,
	QuoteWeekFeildAccount,
	QuoteWeekFeildDateBegin,
	QuoteWeekFeildDateEnd,
	QuoteWeekFeildWeekOfYear,
	QuoteWeekFeildCreateTimestamp,
	QuoteWeekFeildModifyTimestamp,
}

// QuoteWeek quote week
type QuoteWeek struct {
	ID              int64     `json:"id"`
	Code            string    `json:"code"`
	Open            float64   `json:"open"`
	Close           float64   `json:"close"`
	High            float64   `json:"high"`
	Low             float64   `json:"low"`
	Volume          int64     `json:"volume"`
	Account         float64   `json:"account"`
	DateBegin       time.Time `json:"date_begin"`
	DateEnd         time.Time `json:"date_end"`
	WeekOfYear      int       `json:"week_of_year"`
	CreateTimestamp time.Time `json:"create_timestamp"`
	ModifyTimestamp time.Time `json:"modify_timestamp"`
}

func (q *QuoteWeek) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(q)
	return string(buf)
}
