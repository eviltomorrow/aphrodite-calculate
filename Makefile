build: 
	go build -o aphrodite-calculate.runtime -ldflags "-s -w" startup.go