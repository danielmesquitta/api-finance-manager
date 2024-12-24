package root

//go:generate sqlc generate
//go:generate copy-sqlc-params --input ./internal/provider/db/sqlc --output ./internal/provider/repo
//go:generate prisma-to-go entities --schema ./sql/schema.prisma --output ./internal/domain/entity
//go:generate prisma-to-go tables --schema ./sql/schema.prisma --output ./internal/provider/db/query
