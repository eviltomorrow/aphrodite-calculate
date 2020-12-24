package service

import (
	"fmt"

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

	// if len(quotes) != 9 {
	// 	zlog.Warn("[KDJDay]No enough quote data", zap.Int("total-count", len(quotes)), zap.String("code", code), zap.String("date", date))
	// 	return nil, nil
	// }

	var firstDayQuote = quotes[0]
	if firstDayQuote.Date.Format("2006-01-02") != date {
		zlog.Warn("[KDJDay]No exist quote data", zap.String("first-date", firstDayQuote.Date.Format("2006-01-02")), zap.String("code", code), zap.String("date", date))
		return nil, nil
	}

	var c = firstDayQuote.Close

	var lows = make([]float64, 0, len(quotes))
	var highs = make([]float64, 0, len(quotes))
	for _, quote := range quotes {
		lows = append(lows, quote.Low)
		highs = append(highs, quote.High)
	}

	var l = tools.CalcalateMinFloat64(lows)
	var h = tools.CalcalateMaxFloat64(highs)
	var rsv = (c - l) / (h - l) * 100.0
	fmt.Println("rsv: ", rsv)
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

	fmt.Println("kpre: ", kpre)
	fmt.Println("dpre: ", dpre)
	var k = float64(2)/float64(3)*kpre + float64(1)/float64(3)*rsv
	fmt.Println("k: ", k)
	// time.Sleep(1 * time.Second)
	var d = float64(2)/float64(3)*dpre + float64(1)/float64(3)*k
	fmt.Println("d: ", d)
	var j = 3*k - 2*d
	fmt.Println("j: ", j)

	return &model.KDJDay{
		Code:      code,
		K:         tools.Trunc2(k),
		D:         tools.Trunc2(d),
		J:         tools.Trunc2(j),
		Date:      firstDayQuote.Date,
		DayOfYear: firstDayQuote.DayOfYear,
	}, nil
}

// CalculateKDJDay calculate boll day
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
	return 0, nil
}

func kdj(c float64, high []float64, low []float64, kpre, dpre float64) (float64, float64, float64) {
	return 0, 0, 0
}
