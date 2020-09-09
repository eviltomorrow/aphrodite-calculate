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

	log.Printf("model init function\r\n")
}
