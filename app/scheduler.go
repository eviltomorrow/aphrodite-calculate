package app

import (
	"go.uber.org/zap"

	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-calculate/service"
)

//
var (
	DefaultCronSpec               = "01 00 * * MON,TUE,WED,THU,FRI"
	SignalCH        chan struct{} = make(chan struct{}, 64)
)

func runjob() {
	go func() {
		for range SignalCH {
			// 获取 task
			records, err := service.PollUncompletedTaskRecord(false)
			if err != nil {
				zlog.Error("Poll uncompleted task record failure", zap.Error(err))
				continue
			}

			for _, record := range records {
				if record.Completed {
					continue
				}

				switch record.Method {
				case service.SyncQuoteDay:
					if err = service.SyncQuoteDayFromMongoDBToMySQL(record.Date); err != nil {
						zlog.Error("Sync quote day failure", zap.Int64("id", record.ID), zap.String("date", record.Date), zap.Error(err))
					}

				case service.SyncQuoteWeek:
					if err = service.SyncQuoteWeekFromMongoDBToMySQL(record.Date); err != nil {
						zlog.Error("Sync quote week failure", zap.Int64("id", record.ID), zap.String("date", record.Date), zap.Error(err))
					}

				default:
					zlog.Warn("Not support method", zap.String("method", record.Method))
				}
			}
		}
	}()
}

// // Init scheduler
// func initScheduler() error {
// 	var c = cron.New()

// 	if _, err := c.AddFunc(DefaultCronSpec, syncCronFunc); err != nil {
// 		zlog.Fatal("Cron add func failure", zap.Error(err))
// 	}
// 	c.Start()
// 	return nil
// }

// // 同步日数据：每周 2(5)，3(1)，4(2)，5(3)，1(4)
// func syncCronFunc() {
// 	var syncDate string
// 	now := time.Now()

// 	switch now.Weekday() {
// 	case time.Monday, time.Tuesday:
// 		syncDate = now.AddDate(0, 0, -4).Format("2006-01-02")

// 	case time.Wednesday, time.Thursday, time.Friday:
// 		syncDate = now.AddDate(0, 0, -2).Format("2006-01-02")

// 	default:
// 		zlog.Error("Panic: Invalid sync date", zap.Time("now", now))
// 		return
// 	}

// }
