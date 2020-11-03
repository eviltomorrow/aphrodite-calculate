package model

import (
	"time"

	jsoniter "github.com/json-iterator/go"
)

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
	YearWeek        int       `json:"week_of_year"`
	CreateTimestamp time.Time `json:"create_timestamp"`
	ModifyTimestamp time.Time `json:"modify_timestamp"`
}

func (q *QuoteWeek) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(q)
	return string(buf)
}
