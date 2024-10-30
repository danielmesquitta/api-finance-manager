.PHONY: default install update run clear generate build lint migrate seed new_entity

include .env
schema=sql/schema.prisma

default: run

install:
	@go mod download && go install github.com/air-verse/air@latest
update:
	@go mod tidy && go get -u ./...
run:
	@air -c .air.toml
clear:
	@find ./tmp -mindepth 1 ! -name '.gitkeep' -delete
generate:
	@go generate ./...
build:
	@GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ./tmp/restapi ./cmd/restapi
lint:
	@golangci-lint run && nilaway ./...
create_migration:
	@prisma-client-go migrate dev --schema=$(schema)
migrate:
	@prisma-client-go migrate deploy --schema=$(schema)
