package app

import (
	"time"

	"go.uber.org/zap"

	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/eviltomorrow/aphrodite-calculate/service"
	"github.com/robfig/cron"
)

//
var (
	DefaultCronSpec               = "10 23 * * MON,TUE,WED,THU,FRI"
	JobChan         chan struct{} = make(chan struct{}, 64)
	BeginDate                     = "2020-08-31"
)

func initdb() {
	db.BuildMongoDB()
	db.BuildMySQL()
}

func initTaskRecord() {
	begin, err := time.ParseInLocation("2006-01-02", BeginDate, time.Local)
	if err != nil {
		zlog.Fatal("Init: Parse begin date failure", zap.Error(err))
	}

	err = service.BuildTaskRecord(begin, calculateRecordDate())
	if err != nil {
		zlog.Fatal("Init: Build task record failure", zap.Error(err))
	}

}

func initStockList() {
	_, err := service.SyncStockAllFromMongoDBToMySQL()
	if err != nil {
		zlog.Fatal("Init: Sync stock list failure", zap.Error(err))
	}
}

func initCrontab() {
	var c = cron.New()
	_, err := c.AddFunc(DefaultCronSpec, func() {
		var today = calculateRecordDate()

		err := service.BuildTaskRecord(today, today)
		if err != nil {
			zlog.Error("Build task record failure", zap.Time("date", today), zap.Error(err))
		}

		JobChan <- struct{}{}
	})
	if err != nil {
		zlog.Fatal("Init: Build crontab func failure", zap.Error(err))
	}
	c.Start()
}

func runjob() {
	go func() {
		for range JobChan {
			// 获取 task
			records, err := service.PollTaskRecord(false)
			if err != nil {
				zlog.Error("Poll uncompleted task record failure", zap.Error(err))
				continue
			}

			for _, record := range records {
				if record.Completed {
					continue
				}
				record.NumOfTimes = record.NumOfTimes + 1

				var count int64
				switch record.Method {
				case service.SyncQuoteDay:
					if count, err = service.SyncQuoteDayFromMongoDBToMySQL(record.Date); err != nil {
						zlog.Error("Sync quote day failure", zap.Int64("id", record.ID), zap.String("date", record.Date), zap.Error(err))
					}
				case service.SyncQuoteWeek:
					if count, err = service.SyncQuoteWeekFromMongoDBToMySQL(record.Date); err != nil {
						zlog.Error("Sync quote week failure", zap.Int64("id", record.ID), zap.String("date", record.Date), zap.Error(err))
					}

				default:
					zlog.Warn("Not support method", zap.String("method", record.Method))
				}
				if count != 0 {
					record.Completed = true
				}
				if err := service.ArchiveTaskRecord(record); err != nil {
					zlog.Error("Archive task record failure", zap.Int64("id", record.ID), zap.Error(err))
				}
			}

		}
	}()

	JobChan <- struct{}{}
}

func calculateRecordDate() time.Time {
	var now = time.Now()
	var point = time.Date(now.Year(), now.Month(), now.Day(), 23, 10, 0, 0, time.Local)

	switch now.Weekday() {
	case time.Monday:
		if now.After(point) {
			return now.AddDate(0, 0, -3) // 星期五
		}
		return now.AddDate(0, 0, -4) // 星期四

	case time.Tuesday:
		if now.After(point) {
			return now.AddDate(0, 0, -1) // 星期一
		}
		return now.AddDate(0, 0, -4) // 星期五

	case time.Wednesday, time.Thursday, time.Friday:
		if now.After(point) {
			return now.AddDate(0, 0, -1) // 星期二， 星期三， 星期四
		}
		return now.AddDate(0, 0, -2) // 星期一， 星期二， 星期三

	case time.Saturday:
		return now.AddDate(0, 0, -2) // 星期四

	default:
		return now.AddDate(0, 0, -3) // 星期四
	}
}
