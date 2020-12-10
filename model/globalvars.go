package model

import (
	"time"
)

// common
var (
	timeout       = 10 * time.Second
	SelectTimeout = timeout
	InsertTimeout = timeout
	DeleteTimeout = timeout
	UpdateTimeout = timeout
)

// mongodb
var (
	MongodbDatabaseName = "aphrodite"
	CollectionNameQuote = "quote_%s"
	CollectionNameStock = "stock"
)
