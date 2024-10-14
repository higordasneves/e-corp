.PHONY: install-linters
install-linters:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint:
	@echo "==> Running golang ci"
	$$(go env GOPATH)/bin/golangci-lint run --timeout=300s -c ./.golangci.yml ./...

.PHONY: swagger
swagger:
	swag init -g ./cmd/main.go -o ./docs/swagger