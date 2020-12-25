package service

import (
	"github.com/eviltomorrow/aphrodite-base/tools"
	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/eviltomorrow/aphrodite-calculate/model"
	"go.uber.org/zap"
)

func buildKDJDay(code string, date string) (*model.KDJDay, error) {
	quotes, err := model.SelectQuoteDayByCodeDateLatest(db.MySQL, code, date, 9)
	if err != nil {
		return nil, err
	}

	if len(quotes) == 0 {
		zlog.Warn("[KDJDay]No enough quote data", zap.Int("total-count", len(quotes)), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var firstDayQuote = quotes[0]
	if firstDayQuote.Date.Format("2006-01-02") != date {
		zlog.Warn("[KDJDay]No exist quote data", zap.String("first-date", firstDayQuote.Date.Format("2006-01-02")), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var c = firstDayQuote.Close

	var highs = make([]float64, 0, len(quotes))
	var lows = make([]float64, 0, len(quotes))
	for _, quote := range quotes {
		highs = append(highs, quote.High)
		lows = append(lows, quote.Low)
	}

	var kpre float64 = 50.0
	var dpre float64 = 50.0
	kdjs, err := model.SelectKDJDayByCodeDateLimit(db.MySQL, code, date, 1)
	if err != nil {
		return nil, err
	}
	if len(kdjs) == 1 {
		kpre = kdjs[0].K
		dpre = kdjs[0].D
	}

	k, d, j := kdj(c, highs, lows, kpre, dpre)

	return &model.KDJDay{
		Code:      code,
		K:         k,
		D:         d,
		J:         j,
		Date:      firstDayQuote.Date,
		DayOfYear: firstDayQuote.DayOfYear,
	}, nil
}

func buildKDJWeek(code string, date string) (*model.KDJWeek, error) {
	quotes, err := model.SelectQuoteWeekByCodeDateLatest(db.MySQL, code, date, 9)
	if err != nil {
		return nil, err
	}

	if len(quotes) == 0 {
		zlog.Warn("[KDJDay]No enough quote data", zap.Int("total-count", len(quotes)), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var firstWeekQuote = quotes[0]
	if firstWeekQuote.DateEnd.Format("2006-01-02") != date {
		zlog.Warn("[KDJDAY]No exist quote data", zap.String("first-date", firstWeekQuote.DateEnd.Format("2006-01-02")), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var c = firstWeekQuote.Close

	var highs = make([]float64, 0, len(quotes))
	var lows = make([]float64, 0, len(quotes))
	for _, quote := range quotes {
		highs = append(highs, quote.High)
		lows = append(lows, quote.Low)
	}

	var kpre float64 = 50.0
	var dpre float64 = 50.0
	kdjs, err := model.SelectKDJWeekByCodeDateLimit(db.MySQL, code, date, 1)
	if err != nil {
		return nil, err
	}
	if len(kdjs) == 1 {
		kpre = kdjs[0].K
		dpre = kdjs[0].D
	}

	k, d, j := kdj(c, highs, lows, kpre, dpre)

	return &model.KDJWeek{
		Code:       code,
		K:          k,
		D:          d,
		J:          j,
		Date:       firstWeekQuote.DateEnd,
		WeekOfYear: firstWeekQuote.WeekOfYear,
	}, nil
}

// CalculateKDJDay calculate kdj day
func CalculateKDJDay(date string) (int64, error) {
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

		var kdjs = make([]*model.KDJDay, 0, limit)
		var codes = make([]string, 0, limit)

		for _, stock := range stocks {
			codes = append(codes, stock.Code)

			kdj, err := buildKDJDay(stock.Code, date)
			if err != nil {
				return 0, err
			}
			if kdj != nil {
				kdjs = append(kdjs, kdj)
			}
		}

		if len(kdjs) == 0 {
			continue
		}

		tx, err := db.MySQL.Begin()
		if err != nil {
			return 0, err
		}
		_, err = model.DeleteKDJDayByCodesDate(tx, codes, date)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		affected, err := model.InsertKDJDayMany(tx, kdjs)
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

// CalculateKDJWeek calculate kdj week
func CalculateKDJWeek(date string) (int64, error) {
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

		var kdjs = make([]*model.KDJWeek, 0, limit)
		var codes = make([]string, 0, limit)

		for _, stock := range stocks {
			codes = append(codes, stock.Code)

			kdj, err := buildKDJWeek(stock.Code, date)
			if err != nil {
				return 0, err
			}
			if kdj != nil {
				kdjs = append(kdjs, kdj)
			}
		}

		if len(kdjs) == 0 {
			continue
		}

		tx, err := db.MySQL.Begin()
		if err != nil {
			return 0, err
		}
		_, err = model.DeleteKDJWeekByCodesDate(tx, codes, date)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		affected, err := model.InsertKDJWeekMany(tx, kdjs)
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

func kdj(c float64, high []float64, low []float64, kpre, dpre float64) (float64, float64, float64) {
	if len(high) != len(low) {
		return 0, 0, 0
	}
	var h = tools.CalcalateMaxFloat64(high)
	var l = tools.CalcalateMinFloat64(low)
	var n = float64(len(high))
	var m = 1.0

	var rsv = (c - l) / (h - l) * 100
	var k = (m*rsv + (n-m)*kpre) / n
	var d = (m*k + (n-m)*dpre) / n
	var j = 3.0*k - 2.0*d

	return tools.Trunc2(k), tools.Trunc2(d), tools.Trunc2(j)
}
