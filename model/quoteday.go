package model

import (
	"time"

	jsoniter "github.com/json-iterator/go"
)

// QuoteDay quote day
type QuoteDay struct {
	ID              int64     `json:"id"`
	Code            string    `json:"code"`
	Open            float64   `json:"open"`
	Close           float64   `json:"close"`
	High            float64   `json:"high"`
	Low             float64   `json:"low"`
	Volume          int64     `json:"volume"`
	Account         float64   `json:"account"`
	Date            time.Time `json:"date"`
	YearDay         int       `json:"year_day"`
	CreateTimestamp time.Time `json:"create_timestamp"`
	ModifyTimestamp time.Time `json:"modify_timestamp"`
}

func (q *QuoteDay) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(q)
	return string(buf)
}
