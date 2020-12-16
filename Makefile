build: 
	go build -o cmd/aphrodite-calculate.runtime -ldflags "-s -w" startup.go