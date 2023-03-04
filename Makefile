BIN_NAME = rest_api

build:
	@go build -o ./bin/$(BIN_NAME) ./cmd/main.go
run: build
	@./bin/$(BIN_NAME)