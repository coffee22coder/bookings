test-unit:
	go test ./internal/...

test-integration:
	go test -v -tags=integration ./internal/...

test: test-unit test-integration

build:
	go build -o bin/api .