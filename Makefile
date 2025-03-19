include .env
schema=./sql/schema.prisma

ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

define migrate_sequence
	$(MAKE) zip_migrations
	prisma-client-go migrate dev --schema=$(schema) --skip-generate
	$(MAKE) triggers
	$(MAKE) migrate
	$(MAKE) unzip_migrations
endef

define deploy_migrations_sequence
	$(MAKE) zip_migrations
	prisma-client-go migrate deploy --schema=$(schema)
	$(MAKE) unzip_migrations
endef

define reset_db_sequence
	$(MAKE) zip_migrations
	prisma-client-go migrate reset --schema=$(schema) --skip-generate
	$(MAKE) unzip_migrations
endef

.PHONY: default
default: run

.PHONY: install
install:
	@go mod download && ./bin/install.sh

.PHONY: update
update:
	@go mod tidy && go get -u ./...

.PHONY: run
run:
	@air -c .air.toml

.PHONY: clear
clear:
	@find ./tmp -mindepth 1 ! -name '.gitkeep' -delete

.PHONY: generate
generate:
	@go generate ./...

.PHONY: docs
docs:
	@swag init -g ./cmd/restapi/main.go -o ./docs -q && swag2op init -g cmd/restapi/main.go -q --openapiOutputDir ./tmp && mv ./tmp/swagger.json ./docs/openapi.json && mv ./tmp/swagger.yaml ./docs/openapi.yaml

.PHONY: build
build:
	@GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ./tmp/restapi ./cmd/restapi

.PHONY: lint
lint:
	@golangci-lint run && nilaway ./...

.PHONY: lint-fix
lint-fix:
	@golangci-lint run --fix && golines **/*.go -w -m 80

.PHONY: zip_migrations
zip_migrations:
	@prisma-go-tools zip --schema=$(schema) || true

.PHONY: unzip_migrations
unzip_migrations:
	@prisma-go-tools unzip --schema=$(schema)

.PHONY: migrate
migrate:
	@$(migrate_sequence)

.PHONY: create_migration
create_migration:
	@./bin/create_migration.sh sql/migrations $(ARGS)

.PHONY: create_testdata
create_testdata:
	@./bin/create_migration.sh sql/testdata $(ARGS)

.PHONY: deploy_migrations
deploy_migrations:
	@$(deploy_migrations_sequence)

.PHONY: reset_db
reset_db:
	@$(reset_db_sequence)

.PHONY: studio
studio:
	@npx prisma studio --schema=$(schema)

.PHONY: unit-test
unit-test:
	@ENVIRONMENT=test go test -cover -coverprofile=tmp/coverage.out ./internal/domain/usecase/... ./internal/pkg/... -timeout 5s

.PHONY: integration-test
integration-test:
	@ENVIRONMENT=test go test ./test/integration/... -timeout 30s

.PHONY: test
test: unit-test integration-test
	@true

.PHONY: coverage
coverage:
	@go tool cover -html=tmp/coverage.out

.PHONY: triggers
triggers:
	@prisma-go-tools triggers --schema=$(schema)

.PHONY: seed
seed:
	@go run cmd/seed/main.go

%::
	@true
