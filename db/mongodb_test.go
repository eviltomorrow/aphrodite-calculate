package db

import "testing"

func TestBuildMongoDB(t *testing.T) {
	MongodbDSN = "mongodb://localhost:27017"
	MongodbMaxOpen = 10
	MongodbMinOpen = 5
	BuildMongoDB()
}
