package service

import (
	"math"

	"github.com/eviltomorrow/aphrodite-base/tools"

	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-base/ztime"
	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/eviltomorrow/aphrodite-calculate/model"
	"go.uber.org/zap"
)

func buildMADay(code string, date string) (*model.MADay, error) {
	quotes, err := model.SelectQuoteDayByCodeDateLatest(db.MySQL, code, date, 30)
	if err != nil {
		return nil, err
	}

	if len(quotes) == 0 || len(quotes) < 5 {
		zlog.Warn("[MADay]No enough quote data", zap.Int("total-count", len(quotes)), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var firstDayQuote = quotes[0]
	if firstDayQuote.Date.Format("2006-01-02") != date {
		zlog.Warn("[MADAY]No exist quote data", zap.String("first-date", firstDayQuote.Date.Format("2006-01-02")), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var data = make([]float64, 0, len(quotes))
	var m5, m10, m20, m30 float64
	if len(quotes) >= 5 {
		data = data[:0]
		for i := 0; i < 5; i++ {
			data = append(data, quotes[i].Close)
		}
		m5 = ma(data)
	}

	if len(quotes) >= 10 {
		data = data[:0]
		for i := 0; i < 10; i++ {
			data = append(data, quotes[i].Close)
		}
		m10 = ma(data)
	}

	if len(quotes) >= 20 {
		data = data[:0]
		for i := 0; i < 20; i++ {
			data = append(data, quotes[i].Close)
		}
		m20 = ma(data)
	}

	if len(quotes) >= 30 {
		data = data[:0]
		for i := 0; i < 30; i++ {
			data = append(data, quotes[i].Close)
		}
		m30 = ma(data)
	}

	return &model.MADay{
		Code:      code,
		M5:        m5,
		M10:       m10,
		M20:       m20,
		M30:       m30,
		Date:      firstDayQuote.Date,
		DayOfYear: firstDayQuote.DayOfYear,
	}, nil
}

func buildMAWeek(code string, date string) (*model.MAWeek, error) {
	quotes, err := model.SelectQuoteWeekByCodeDateLatest(db.MySQL, code, date, 30)
	if err != nil {
		return nil, err
	}

	if len(quotes) == 0 || len(quotes) < 5 {
		zlog.Warn("[MAWeek]No enough quote data", zap.Int("total-count", len(quotes)), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var firstDayQuote = quotes[0]
	if firstDayQuote.DateEnd.Format("2006-01-02") != date {
		zlog.Warn("[MAWeek]No exist quote data", zap.String("first-date", firstDayQuote.DateEnd.Format("2006-01-02")), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var data = make([]float64, 0, len(quotes))
	var m5, m10, m20, m30 float64
	if len(quotes) >= 5 {
		data = data[:0]
		for i := 0; i < 5; i++ {
			data = append(data, quotes[i].Close)
		}
		m5 = ma(data)
	}

	if len(quotes) >= 10 {
		data = data[:0]
		for i := 0; i < 10; i++ {
			data = append(data, quotes[i].Close)
		}
		m10 = ma(data)
	}

	if len(quotes) >= 20 {
		data = data[:0]
		for i := 0; i < 20; i++ {
			data = append(data, quotes[i].Close)
		}
		m20 = ma(data)
	}

	if len(quotes) >= 30 {
		data = data[:0]
		for i := 0; i < 30; i++ {
			data = append(data, quotes[i].Close)
		}
		m30 = ma(data)
	}

	return &model.MAWeek{
		Code:       code,
		M5:         m5,
		M10:        m10,
		M20:        m20,
		M30:        m30,
		Date:       firstDayQuote.DateEnd,
		WeekOfYear: ztime.YearWeek(firstDayQuote.DateEnd),
	}, nil
}

// CalculateMADay calculate ma day
func CalculateMADay(date string) (int64, error) {
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

		var mas = make([]*model.MADay, 0, limit)
		var codes = make([]string, 0, limit)

		for _, stock := range stocks {
			codes = append(codes, stock.Code)

			ma, err := buildMADay(stock.Code, date)
			if err != nil {
				return 0, err
			}
			if ma != nil {
				mas = append(mas, ma)
			}
		}

		if len(mas) == 0 {
			continue
		}

		tx, err := db.MySQL.Begin()
		if err != nil {
			return 0, err
		}
		_, err = model.DeleteMADayByCodesDate(tx, codes, date)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		affected, err := model.InsertMADayMany(tx, mas)
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
	return 0, nil
}

// CalculateMAWeek calculate ma day
func CalculateMAWeek(date string) (int64, error) {
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

		var mws = make([]*model.MAWeek, 0, limit)
		var codes = make([]string, 0, limit)

		for _, stock := range stocks {
			codes = append(codes, stock.Code)

			mw, err := buildMAWeek(stock.Code, date)
			if err != nil {
				return 0, err
			}
			if mw != nil {
				mws = append(mws, mw)
			}
		}

		if len(mws) == 0 {
			continue
		}

		tx, err := db.MySQL.Begin()
		if err != nil {
			return 0, err
		}
		_, err = model.DeleteMAWeekByCodesDate(tx, codes, date)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		affected, err := model.InsertMAWeekMany(tx, mws)
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
	return 0, nil
}

var n10 = math.Pow10(2)

func ma(close []float64) float64 {
	if len(close) == 0 {
		return 0
	}

	var sum float64
	for _, c := range close {
		sum += c
	}

	return tools.Trunc2(sum / float64(len(close)))
}
