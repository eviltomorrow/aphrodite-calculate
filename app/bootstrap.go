package app

// Startup startup service
func Startup() {
	initCrontab()
	initdb()
	initTaskRecord()
	initStockList()

	runjob()
}
