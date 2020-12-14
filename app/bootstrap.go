package app

import (
	"github.com/eviltomorrow/aphrodite-calculate/db"
)

// StartupService startup service
func StartupService() error {
	db.BuildMongoDB()
	db.BuildMySQL()

	runjob()
	return nil
}
