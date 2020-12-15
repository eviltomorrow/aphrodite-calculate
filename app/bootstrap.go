package app

import (
	"time"

	"github.com/robfig/cron/v3"

	"github.com/eviltomorrow/aphrodite-calculate/db"
)

// StartupService startup service
func StartupService() error {
	db.BuildMongoDB()
	db.BuildMySQL()

	if err := initjob(); err != nil {
		return err
	}

	runjob()

	var c = cron.New()
	_, err := c.AddFunc(DefaultCronSpec, func() {
		DateCH <- time.Now().Format("2006-01-02")
	})
	if err != nil {
		return err
	}
	c.Start()
	return nil
}
