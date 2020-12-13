package app

import (
	"time"

	"github.com/eviltomorrow/aphrodite-calculate/db"
	"github.com/eviltomorrow/aphrodite-calculate/service"
)

//
var (
	BeginDate = "2020-08-31" // begin date
)

// initTaskRecord init task
func initTaskRecord() error {
	beginDate, err := time.ParseInLocation("2006-01-02", BeginDate, time.Local)
	if err != nil {
		return err
	}
	endDate := time.Now()

	return service.BuildTaskRecord(beginDate, endDate)
}

// StartupService startup service
func StartupService() error {
	db.BuildMongoDB()
	db.BuildMySQL()

	if err := initTaskRecord(); err != nil {
		return err
	}

	if err := initScheduler(); err != nil {
		return err
	}

	return nil
}
