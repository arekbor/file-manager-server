BIN_NAME = file_manager_server

build:
	@go build -o ./bin/$(BIN_NAME) ./cmd/main.go
run: build
	@./bin/$(BIN_NAME)