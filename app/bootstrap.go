package app

import (
	"time"

	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-calculate/db"
)

// Startup startup service
func Startup() error {
	db.BuildMongoDB()
	db.BuildMySQL()

	if err := initjob(); err != nil {
		return err
	}
	zlog.Info("Init job complete")

	runjob()

	DateCH <- time.Now().Format("2006-01-02")
	// var c = cron.New()
	// _, err := c.AddFunc(DefaultCronSpec, func() {
	// 	DateCH <- time.Now().Format("2006-01-02")
	// })
	// if err != nil {
	// 	return err
	// }
	// go func() { c.Start() }()
	return nil
}
