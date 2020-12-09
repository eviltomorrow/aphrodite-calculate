package service

import (
	"fmt"

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

// SyncQuoteDayFromMongoDBToMySQL sync quoteday from mongodb to mysql
func SyncQuoteDayFromMongoDBToMySQL() {

}

// SyncQuoteWeekFromMongoDBToMySQL sync quoteweek from mongodb to mysql
func SyncQuoteWeekFromMongoDBToMySQL() {

}
