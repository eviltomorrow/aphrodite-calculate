package model

import (
	"fmt"
	"time"
)

// common
var (
	timeout       = 5 * time.Second
	SelectTimeout = timeout
	InsertTimeout = timeout
	DeleteTimeout = timeout
	UpdateTimeout = timeout
)

// mongodb
var (
	MongodbDatabaseName = "aphrodite"
	CollectionNameQuote = fmt.Sprintf("quote_%s", time.Now().Format("2006"))
	CollectionNameStock = "stock"
)
