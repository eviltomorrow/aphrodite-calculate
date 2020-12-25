package service

import (
	"math"

	"github.com/eviltomorrow/aphrodite-base/tools"
	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/eviltomorrow/aphrodite-calculate/model"

	"go.uber.org/zap"
)

func buildBollDay(code string, date string) (*model.BollDay, error) {
	quotes, err := model.SelectQuoteDayByCodeDateLatest(db.MySQL, code, date, 20)
	if err != nil {
		return nil, err
	}

	if len(quotes) != 20 {
		zlog.Warn("[BollDay]No enough quote data", zap.Int("total-count", len(quotes)), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var firstDayQuote = quotes[0]
	if firstDayQuote.Date.Format("2006-01-02") != date {
		zlog.Warn("[BollDAY]No exist quote data", zap.String("first-date", firstDayQuote.Date.Format("2006-01-02")), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var data = make([]float64, 0, len(quotes))
	for _, quote := range quotes {
		data = append(data, quote.Close)
	}

	up, mb, dn := boll(data)

	return &model.BollDay{
		Code:      code,
		UP:        up,
		MB:        mb,
		DN:        dn,
		Date:      firstDayQuote.Date,
		DayOfYear: firstDayQuote.DayOfYear,
	}, nil
}

func buildBollWeek(code string, date string) (*model.BollWeek, error) {
	quotes, err := model.SelectQuoteWeekByCodeDateLatest(db.MySQL, code, date, 20)
	if err != nil {
		return nil, err
	}

	if len(quotes) != 20 {
		zlog.Warn("[BollDay]No enough quote data", zap.Int("total-count", len(quotes)), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var firstDayQuote = quotes[0]
	if firstDayQuote.DateEnd.Format("2006-01-02") != date {
		zlog.Warn("[BollDAY]No exist quote data", zap.String("first-date", firstDayQuote.DateEnd.Format("2006-01-02")), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var data = make([]float64, 0, len(quotes))
	for _, quote := range quotes {
		data = append(data, quote.Close)
	}

	up, mb, dn := boll(data)

	return &model.BollWeek{
		Code:       code,
		UP:         up,
		MB:         mb,
		DN:         dn,
		Date:       firstDayQuote.DateEnd,
		WeekOfYear: firstDayQuote.WeekOfYear,
	}, nil
}

// CalculateBollDay calculate boll day
func CalculateBollDay(date string) (int64, error) {
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

		var bolls = make([]*model.BollDay, 0, limit)
		var codes = make([]string, 0, limit)

		for _, stock := range stocks {
			codes = append(codes, stock.Code)

			boll, err := buildBollDay(stock.Code, date)
			if err != nil {
				return 0, err
			}
			if boll != nil {
				bolls = append(bolls, boll)
			}
		}

		if len(bolls) == 0 {
			continue
		}

		tx, err := db.MySQL.Begin()
		if err != nil {
			return 0, err
		}
		_, err = model.DeleteBollDayByCodesDate(tx, codes, date)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		affected, err := model.InsertBollDayMany(tx, bolls)
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

// CalculateBollWeek calculate boll week
func CalculateBollWeek(date string) (int64, error) {
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

		var bws = make([]*model.BollWeek, 0, limit)
		var codes = make([]string, 0, limit)

		for _, stock := range stocks {
			codes = append(codes, stock.Code)

			bw, err := buildBollWeek(stock.Code, date)
			if err != nil {
				return 0, err
			}
			if bw != nil {
				bws = append(bws, bw)
			}
		}

		if len(bws) == 0 {
			continue
		}

		tx, err := db.MySQL.Begin()
		if err != nil {
			return 0, err
		}
		_, err = model.DeleteBollWeekByCodesDate(tx, codes, date)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		affected, err := model.InsertBollWeekMany(tx, bws)
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

func boll(close []float64) (float64, float64, float64) {
	if len(close) == 0 || len(close) == 1 {
		return 0, 0, 0
	}

	var sum float64
	for _, c := range close {
		sum += c
	}
	var mb = sum / float64(len(close))

	sum = 0
	for _, c := range close {
		var md = math.Pow((c - mb), 2)
		sum += md
	}

	var sd = math.Sqrt(sum / float64(len(close)-1))
	return tools.Trunc2(mb + 2*sd), tools.Trunc2(mb), tools.Trunc2(mb - 2*sd)
}
