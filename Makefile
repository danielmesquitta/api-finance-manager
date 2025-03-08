.PHONY: default install update run clear generate build lint create_migration migrate reset_db docs test seed

include .env
schema=sql/schema.prisma

define create_migration_sequence
	$(MAKE) zip_migrations
	prisma-client-go migrate dev --schema=$(schema) --skip-generate && prisma-go-tools triggers --schema=$(schema)
	$(MAKE) migrate
	$(MAKE) unzip_migrations
endef

define migrate_sequence
	$(MAKE) zip_migrations
	prisma-client-go migrate deploy --schema=$(schema)
	$(MAKE) unzip_migrations
endef

define reset_db_sequence
	$(MAKE) zip_migrations
	prisma-client-go migrate reset --schema=$(schema) --skip-generate
	$(MAKE) unzip_migrations
endef

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
docs:
	@swag init -g ./cmd/restapi/main.go -o ./docs -q && swag2op init -g cmd/restapi/main.go -q --openapiOutputDir ./tmp && mv ./tmp/swagger.json ./docs/openapi.json && mv ./tmp/swagger.yaml ./docs/openapi.yaml
build:
	@GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ./tmp/restapi ./cmd/restapi
lint:
	@golangci-lint run && nilaway ./...
lint-fix:
	@golangci-lint run --fix && golines **/*.go -w -m 80 && go run cmd/lintfix/main.go
zip_migrations:
	@prisma-go-tools zip --schema ./sql/schema.prisma
unzip_migrations:
	@prisma-go-tools unzip --schema ./sql/schema.prisma
create_migration:
	@$(create_migration_sequence)
migrate:
	@$(migrate_sequence)
reset_db:
	@$(reset_db_sequence)
studio:
	@npx prisma studio --schema=$(schema)
test:
	@ENVIRONMENT=test go test ./test/integration/...
seed:
	@go run cmd/seed/main.go
