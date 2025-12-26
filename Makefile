.PHONY: all install test buf eny prepare build run run-release migrate vulncheck lint

include migration.mk

install:
	@go mod tidy
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/swaggo/swag/cmd/swag@v1.8.7
prepare: install

build:
	go mod download
	go build -o ./src/build ./cmd/executor

run-docker:
	docker-compose up -d

run: doc
	go run ./cmd/executor

test:
	go test ./...

vulncheck:
	govulncheck -tags release ./cmd/executor

lint: doc
	gci write .
	gofumpt -l -w .
	golangci-lint run  -v

doc:
	@cd ./cmd/executor && swag init --parseDependency=true --output "./docs"
