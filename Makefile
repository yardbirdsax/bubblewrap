.PHONY: coverage
coverage:
	go tool cover -func=coverage.out
.PHONY: test
test:
	go test -v -coverprofile coverage.out -race ./...
.PHONY: tidy
tidy:
	go mod tidy