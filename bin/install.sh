#!/bin/bash

packages=(
    "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    "go.uber.org/nilaway/cmd/nilaway@latest"
    "github.com/air-verse/air@latest"
    "github.com/steebchen/prisma-client-go@latest"
    "github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
    "github.com/google/wire/cmd/wire@latest"
)

echo "Installing and updating Go packages..."

for package in "${packages[@]}"; do
    echo "Installing and/or updating $package..."
    go install "$package"
done

echo "All packages have been successfully installed and updated."
