package service

import (
	"fmt"
	"time"

	"github.com/eviltomorrow/aphrodite-base/zlog"
	"go.uber.org/zap"

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
		return 0, "", nil
	}

	tx, err := db.MySQL.Begin()
	if err != nil {
		return 0, "", fmt.Errorf("Get mysql transaction failure, nest error: %v", err)
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
func SyncStockAllFromMongoDBToMySQL() {
	var offset int64 = 0
	var limit int64 = 100
	var lastID string
	var count int64
	for {
		affected, lastID, err := syncStockFromMongoDBToMySQL(offset, limit, lastID)
		if err != nil {
			zlog.Error("Sync stock from mongodb to mysql failure", zap.Int64("offset", offset), zap.Int64("limit", limit), zap.String("lastID", lastID))
		}
		if lastID == "" {
			break
		}
		offset += limit
		count += affected
	}
	zlog.Info("Sync stock complete", zap.Int64("total count", count))
}

func buildQuoteDayFromMongoDBToMySQL(code string, date string) (*model.QuoteDay, error) {
	quotes, err := model.QueryQuoteBaseCurrentCodeLimit2(db.MongoDB, code, date)
	if err != nil {
		return nil, err
	}
	if len(quotes) != 2 {
		zlog.Warn("No enough quote data", zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var firstDayQuote = quotes[0]
	if firstDayQuote.Date != date {
		zlog.Warn("No exist quote data", zap.String("code", code), zap.String("date", date))
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
func SyncQuoteDayFromMongoDBToMySQL(date string) error {
	var offset int64 = 0
	var limit int64 = 50
	for {
		stocks, err := model.SelectStockManyForMySQL(db.MySQL, offset, limit)
		if err != nil {
			return err
		}

		if len(stocks) == 0 {
			break
		}

		var quotes = make([]*model.QuoteDay, 0, limit)
		var codes = make([]string, 0, limit)
		for _, stock := range stocks {
			codes = append(codes, stock.Code)

			quote, err := buildQuoteDayFromMongoDBToMySQL(stock.Code, date)
			if err != nil {
				return err
			}
			if quote == nil {
				continue
			}
			quotes = append(quotes, quote)
		}

		if len(quotes) == 0 {
			continue
		}

		tx, err := db.MySQL.Begin()
		if err != nil {
			return err
		}
		_, err = model.DeleteQuoteDayByCodesDate(tx, codes, date)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = model.InsertQuoteDayMany(tx, quotes)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		offset += limit
	}
	return nil
}

// SyncQuoteWeekFromMongoDBToMySQL sync quoteweek from mongodb to mysql
func SyncQuoteWeekFromMongoDBToMySQL() {

}
