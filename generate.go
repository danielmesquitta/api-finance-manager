package root

//go:generate prisma-go-tools unzip --schema ./sql/schema.prisma
//go:generate sqlc generate
//go:generate copy-sqlc-params --input ./internal/provider/db/sqlc --output ./internal/provider/repo
//go:generate prisma-go-tools entities --schema ./sql/schema.prisma --output ./internal/domain/entity
//go:generate prisma-go-tools tables --schema ./sql/schema.prisma --output ./internal/provider/db/schema
//go:generate wire-config -c internal/config/wire/wire.go -o internal/app/server/wire.go -m github.com/danielmesquitta/api-finance-manager/internal/app/server -e dev,staging,test,prod
