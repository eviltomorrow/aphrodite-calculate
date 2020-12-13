package app

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-calculate/service"
)

//
var (
	DefaultCronSpec = "01 00 * * MON,TUE,WED,THU,FRI"
)

// Init scheduler
func initScheduler() error {
	var c = cron.New()

	if _, err := c.AddFunc(DefaultCronSpec, syncCronFunc); err != nil {
		zlog.Fatal("Cron add func failure", zap.Error(err))
	}
	c.Start()
	return nil
}

// 同步日数据：每周 2(5)，3(1)，4(2)，5(3)，1(4)
func syncCronFunc() {
	var syncDate string
	now := time.Now()

	switch now.Weekday() {
	case time.Monday, time.Tuesday:
		syncDate = now.AddDate(0, 0, -4).Format("2006-01-02")

	case time.Wednesday, time.Thursday, time.Friday:
		syncDate = now.AddDate(0, 0, -2).Format("2006-01-02")

	default:
		zlog.Error("Panic: Invalid sync date", zap.Time("now", now))
		return
	}

	if err := sync(syncDate); err != nil {
		zlog.Error("Sync date failure", zap.Error(err), zap.Time("now", now), zap.String("sync-date", syncDate))
	}
}

func sync(date string) error {
	d, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return err
	}

	service.SyncStockAllFromMongoDBToMySQL()
	if err := service.SyncQuoteDayFromMongoDBToMySQL(date); err != nil {
		return fmt.Errorf("sync quote day failure, nest error: %v, date: %s", err, date)
	}

	// is sync week day
	if d.Weekday() == time.Thursday {
		if err := service.SyncQuoteWeekFromMongoDBToMySQL(date); err != nil {
			return fmt.Errorf("sync quote week failure, nest error: %v, date: %s", err, date)
		}
	}

	return nil
}
