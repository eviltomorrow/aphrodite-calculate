package app

import (
	"time"

	"go.uber.org/zap"

	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-calculate/service"
)

//
var (
	DefaultCronSpec               = "10 23 * * MON,TUE,WED,THU,FRI"
	JobChan         chan struct{} = make(chan struct{}, 64)
	BeginDate                     = "2020-08-31"
)

func initjob() error {
	begin, err := time.ParseInLocation("2006-01-02", BeginDate, time.Local)
	if err != nil {
		return err
	}

	err = service.BuildTaskRecord(begin, time.Now())
	if err != nil {
		return err
	}

	_, err = service.SyncStockAllFromMongoDBToMySQL()
	if err != nil {
		return err
	}

	return nil
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

			var ids = make([]int64, 0, 128)
			for _, record := range records {
				if record.Completed {
					continue
				}

				zlog.Info("Run job", zap.String("date", record.Date))
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
					ids = append(ids, record.ID)
				}
			}
			if err := service.ArchiveTaskRecord(ids); err != nil {
				zlog.Error("Archive task record failure", zap.Any("ids", ids))
			}
		}
	}()
}
