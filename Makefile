build:
	@go build -o bin/distributed-webscraper cmd/main.go

run: build
	@bin/distributed-webscraper

test:
	@go test -v ./..
