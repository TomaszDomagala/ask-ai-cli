.PHONY: build
build:
	go build -o aai

.PHONY: test
test:
	go test ./...
