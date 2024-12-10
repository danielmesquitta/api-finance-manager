.PHONY: default install update run clear generate build lint create_migration migrate

include .env
schema=sql/schema.prisma

default: run

install:
	@go mod download && ./bin/install.sh
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
	@golangci-lint run && nilaway ./... && golines **/*.go -w -m 80
create_migration:
	@prisma-client-go migrate dev --schema=$(schema) --skip-generate
migrate:
	@prisma-client-go migrate deploy --schema=$(schema)
