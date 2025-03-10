target = "monkey"

test_cover:
	go test ./... -v -cover

broswer_test:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out

clean:
	rm -f ${target}

build: clean test_cover
	go build -o ${target} .


all: build

.PHONY: clean build test_cover broswer_test


