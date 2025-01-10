build:
	@go build -o bin/distributed-webscraper main.go

run: build
	@bin/distributed-webscraper

test:
	@go test -v ./..
