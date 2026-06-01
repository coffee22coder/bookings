BIN := bin/api
ifeq ($(OS),Windows_NT)
	BIN := bin/api.exe
endif

test-unit:
	go test ./internal/...

test-integration:
	go test -v -tags=integration ./internal/...

test: test-unit test-integration

build:
	go build -o $(BIN) ./cmd/api