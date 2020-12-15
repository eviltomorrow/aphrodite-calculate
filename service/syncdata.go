package service

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/eviltomorrow/aphrodite-base/tools"
	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-base/ztime"
	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/eviltomorrow/aphrodite-calculate/model"
)

func syncStockFromMongoDBToMySQL(offset, limit int64, lastID string) (int64, string, error) {
	stocks, err := model.SelectStockManyForMongoDB(db.MongoDB, offset, limit, lastID)
	if err != nil {
		return 0, "", fmt.Errorf("Query stock from mongodb failure, nest error: %v", err)
	}

	if len(stocks) == 0 {
		return 0, "", nil
	}
	lastID = stocks[len(stocks)-1].ObjectID

	var codes = make([]string, 0, len(stocks))
	for _, stock := range stocks {
		codes = append(codes, stock.Code)
	}

	data, err := model.SelectStockManyByCodesForMySQL(db.MySQL, codes)
	if err != nil {
		return 0, lastID, fmt.Errorf("Query stock from mysql failure, nest error: %v", err)
	}

	var updateCache = make([]*model.Stock, 0, len(stocks)/2)
	var insertCache = make([]*model.Stock, 0, len(stocks)/2)
loop:
	for _, stock := range stocks {
		for _, tmp := range data {
			if stock.Code == tmp.Code {
				if stock.Name != tmp.Name || stock.Source != tmp.Source {
					updateCache = append(updateCache, stock)
				}
				continue loop
			}
		}
		insertCache = append(insertCache, stock)
	}

	if len(updateCache) == 0 && len(insertCache) == 0 {
		return 0, lastID, nil
	}

	tx, err := db.MySQL.Begin()
	if err != nil {
		return 0, lastID, fmt.Errorf("Get mysql transaction failure, nest error: %v", err)
	}

	var count int64
	for _, stock := range updateCache {
		affected, err := model.UpdateStockByCodeForMySQL(tx, stock.Code, stock)
		if err != nil {
			tx.Rollback()
			return 0, lastID, fmt.Errorf("Update stock for mysql failure, nest error: %v", err)
		}
		count += affected
	}

	affected, err := model.InsertStockManyForMySQL(tx, insertCache)
	if err != nil {
		tx.Rollback()
		return 0, lastID, fmt.Errorf("Insert stock for mysql failure, nest error: %v", err)
	}
	if err := tx.Commit(); err != nil {
		return 0, lastID, fmt.Errorf("Commit transaction failure, nest error: %v", err)
	}

	if len(stocks) < int(limit) {
		lastID = ""
	}
	count += affected

	return count, lastID, nil
}

// SyncStockAllFromMongoDBToMySQL sync stock from mongodb to mysql
func SyncStockAllFromMongoDBToMySQL() (int64, error) {
	var offset int64 = 0
	var limit int64 = 100
	var lastID string
	var count int64
	for {
		affected, lastID, err := syncStockFromMongoDBToMySQL(offset, limit, lastID)
		if err != nil {
			return 0, fmt.Errorf("Sync stock failure, nest error: %v, offset: %d, limit: %d, lastID: %s", err, offset, limit, lastID)
		}
		if lastID == "" {
			break
		}
		offset += limit
		count += affected
	}
	return count, nil
}

func buildQuoteDayFromMongoDBToMySQL(code string, date string) (*model.QuoteDay, error) {
	quotes, err := model.SelectQuoteBaseCurrentCodeLimit2(db.MongoDB, code, date)
	if err != nil {
		return nil, err
	}

	if len(quotes) != 2 {
		zlog.Warn("[QuoteBase]No enough quote data", zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var firstDayQuote = quotes[0]
	if firstDayQuote.Date != date {
		zlog.Warn("[QuoteBase]No exist quote data", zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	d, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return nil, err
	}
	var secondDayQuote = quotes[1]
	var quoteDay = &model.QuoteDay{
		Code:      code,
		Open:      firstDayQuote.Open,
		Close:     secondDayQuote.YesterdayClosed,
		High:      firstDayQuote.High,
		Low:       firstDayQuote.Low,
		Volume:    firstDayQuote.Volume,
		Account:   firstDayQuote.Account,
		Date:      d,
		DayOfYear: d.YearDay(),
	}
	return quoteDay, nil
}

// SyncQuoteDayFromMongoDBToMySQL sync quoteday from mongodb to mysql
func SyncQuoteDayFromMongoDBToMySQL(date string) (int64, error) {
	var offset int64 = 0
	var limit int64 = 50
	var count int64
	for {
		stocks, err := model.SelectStockManyForMySQL(db.MySQL, offset, limit)
		if err != nil {
			return 0, err
		}

		offset += limit
		if len(stocks) == 0 {
			break
		}

		var quotes = make([]*model.QuoteDay, 0, limit)
		var codes = make([]string, 0, limit)
		for _, stock := range stocks {
			codes = append(codes, stock.Code)

			quote, err := buildQuoteDayFromMongoDBToMySQL(stock.Code, date)
			if err != nil {
				return 0, err
			}
			if quote != nil {
				quotes = append(quotes, quote)
			}

		}

		if len(quotes) == 0 {
			continue
		}

		tx, err := db.MySQL.Begin()
		if err != nil {
			return 0, err
		}
		_, err = model.DeleteQuoteDayByCodesDate(tx, codes, date)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		affected, err := model.InsertQuoteDayMany(tx, quotes)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		if err = tx.Commit(); err != nil {
			tx.Rollback()
			return 0, err
		}

		count += affected
	}
	return count, nil
}

func buildQuoteWeekFromQuoteDay(code string, begin, end time.Time) (*model.QuoteWeek, error) {
	quotes, err := model.SelectQuoteDayByCodeDate(db.MySQL, code, begin.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	if len(quotes) == 0 {
		zlog.Warn("[QuoteDay]No exist quote data", zap.String("code", code), zap.Time("begin", begin), zap.Time("end", end))
		return nil, nil
	}

	var quoteWeek = &model.QuoteWeek{
		Code:       code,
		Open:       quotes[0].Open,
		Close:      quotes[len(quotes)-1].Close,
		DateBegin:  begin,
		DateEnd:    end,
		WeekOfYear: ztime.YearWeek(begin),
	}
	var numFloat64 = make([]float64, 0, len(quotes))
	var numInt64 = make([]int64, 0, len(quotes))

	// cal high
	numFloat64 = numFloat64[:0]
	for _, quote := range quotes {
		numFloat64 = append(numFloat64, quote.High)
	}
	quoteWeek.High = tools.CalcalateMaxFloat64(numFloat64)

	// cal low
	numFloat64 = numFloat64[:0]
	for _, quote := range quotes {
		numFloat64 = append(numFloat64, quote.Low)
	}
	quoteWeek.Low = tools.CalcalateMinFloat64(numFloat64)

	// cal volume
	numInt64 = numInt64[:0]
	for _, quote := range quotes {
		numInt64 = append(numInt64, quote.Volume)
	}
	quoteWeek.Volume = tools.CalculateSumInt64(numInt64)

	// cal account
	numFloat64 = numFloat64[:0]
	for _, quote := range quotes {
		numFloat64 = append(numFloat64, quote.Account)
	}
	quoteWeek.Account = tools.CalculateSumFloat64(numFloat64)

	return quoteWeek, nil
}

// SyncQuoteWeekFromMongoDBToMySQL sync quoteweek from mongodb to mysql
func SyncQuoteWeekFromMongoDBToMySQL(date string) (int64, error) {
	var offset int64 = 0
	var limit int64 = 50
	var count int64
	end, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return 0, err
	}

	begin := end.AddDate(0, 0, -5)

	for {
		stocks, err := model.SelectStockManyForMySQL(db.MySQL, offset, limit)
		if err != nil {
			return 0, err
		}

		if len(stocks) == 0 {
			break
		}

		var quotes = make([]*model.QuoteWeek, 0, limit)
		var codes = make([]string, 0, limit)
		for _, stock := range stocks {
			codes = append(codes, stock.Code)
			quote, err := buildQuoteWeekFromQuoteDay(stock.Code, begin, end)
			if err != nil {
				return 0, err
			}
			if quote != nil {
				quotes = append(quotes, quote)
			}
		}

		if len(quotes) == 0 {
			continue
		}

		tx, err := db.MySQL.Begin()
		if err != nil {
			return 0, err
		}
		_, err = model.DeleteQuoteWeekByCodesDate(tx, codes, end.Format("2006-01-02"))
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		affected, err := model.InsertQuoteWeekMany(tx, quotes)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		if err = tx.Commit(); err != nil {
			tx.Rollback()
			return 0, err
		}

		offset += limit
		count += affected
	}
	return count, nil
}
