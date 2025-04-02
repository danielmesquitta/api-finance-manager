#!/bin/bash

packages=(
    "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2"
    "go.uber.org/nilaway/cmd/nilaway@latest"
    "github.com/segmentio/golines@latest"
    "github.com/air-verse/air@latest"
    "github.com/steebchen/prisma-client-go@latest"
    "github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
    "github.com/google/wire/cmd/wire@latest"
    "github.com/danielmesquitta/prisma-go-tools@latest"
    "github.com/danielmesquitta/copy-sqlc-params@latest"
    "github.com/danielmesquitta/wire-config@latest"
    "github.com/swaggo/swag/cmd/swag@latest"
    "github.com/zxmfke/swagger2openapi3/cmd/swag2op@latest"
)

echo "Installing and updating Go packages..."

for package in "${packages[@]}"; do
    echo "$package..."
    go install "$package"
done

echo "All packages have been successfully installed and updated."
