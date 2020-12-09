package service

import (
	"log"
	"os"

	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-calculate/db"
)

func init() {
	global, prop, err := zlog.InitLogger(&zlog.Config{
		Level:            "debug",
		Format:           "text",
		DisableTimestamp: false,
		File: zlog.FileLogConfig{
			Filename: "/tmp/aphrodite-calculate/data.log",
			MaxSize:  30,
		},
		DisableStacktrace: true,
	})
	if err != nil {
		log.Printf("Fatal: Setup log config failure, nest error: %v", err)
		os.Exit(1)
	}
	zlog.ReplaceGlobals(global, prop)

	db.MongodbDSN = "mongodb://localhost:27017"
	db.MongodbMaxOpen = 10
	db.MongodbMinOpen = 5
	db.BuildMongoDB()

	db.MySQLDSN = "root:root@tcp(localhost:3306)/aphrodite?charset=utf8mb4&parseTime=true&loc=Local"
	db.MySQLMaxOpen = 10
	db.MySQLMinOpen = 5
	db.BuildMySQL()
	log.Printf("model init function\r\n")
}
