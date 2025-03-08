package root

//go:generate prisma-go-tools unzip --schema ./sql/schema.prisma
//go:generate sqlc generate
//go:generate copy-sqlc-params --input ./internal/provider/db/sqlc --output ./internal/provider/repo
//go:generate prisma-go-tools entities --schema ./sql/schema.prisma --output ./internal/domain/entity
//go:generate prisma-go-tools tables --schema ./sql/schema.prisma --output ./internal/provider/db
//go:generate go run ./cmd/wire/main.go
