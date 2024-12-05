package main

//go:generate sqlc generate
//go:generate wire ./...
//go:generate prisma-to-go --schema ./sql/schema.prisma --output ./internal/domain/entity
