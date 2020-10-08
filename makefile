test-cover:
	go test -v -coverprofile .testCoverage.txt
	go tool cover -html=.testCoverage.txt -o coverage.html

run:
	go run .

build:
	go build