.PHONY: coverage
coverage:
	go tool cover -func=coverage.out
.PHONY: test
test:
	go test -v -coverprofile coverage.out ./...
.PHONY: tidy
tidy:
	go mod tidy
.PHONY: fmtcheck
fmtcheck:
	@test -z "$$(go fmt ./...)"
.PHONY: fmt
	@go fmt ./...