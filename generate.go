package root

//go:generate sqlc generate
//go:generate copy-sqlc-params --input ./internal/provider/db/sqlc --output ./internal/provider/repo
//go:generate prisma-to-go --schema ./sql/schema.prisma --output ./internal/domain/entity
//go:generate swag init -g ./cmd/restapi/main.go -o ./docs
