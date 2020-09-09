package db

import "testing"

func TestBuildMySQL(t *testing.T) {
	MySQLDSN = "root:root@tcp(localhost:3306)/toffee?charset=utf8"
	MySQLMaxOpen = 10
	MySQLMinOpen = 5
	BuildMySQL()
}
