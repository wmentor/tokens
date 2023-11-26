.PHONY: all test

all: test

test:
	go clean -testcache
	go test ./... -cover
