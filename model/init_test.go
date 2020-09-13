package model

import (
	"log"

	"github.com/eviltomorrow/aphrodite-calculate/db"
)

func init() {
	db.MongodbDSN = "mongodb://localhost:27017"
	db.MongodbMaxOpen = 10
	db.MongodbMinOpen = 5
	db.BuildMongoDB()

	db.MySQLDSN = "root:root@tcp(localhost:3306)/aphrodite?charset=utf8"
	db.MySQLMaxOpen = 10
	db.MySQLMinOpen = 5
	db.BuildMySQL()
	log.Printf("model init function\r\n")
}
