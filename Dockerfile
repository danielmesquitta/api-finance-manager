FROM golang:1.23.4-alpine as builder
WORKDIR /app
COPY . .
# Migrate
RUN go install github.com/steebchen/prisma-client-go@latest
RUN prisma-client-go migrate deploy --schema=sql/schema.prisma
# Build
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ./tmp/restapi ./cmd/restapi

FROM alpine:latest
COPY --from=builder /app/tmp/restapi .
RUN apk upgrade --no-cache
RUN apk add --no-cache ca-certificates
EXPOSE 8080
CMD ["./restapi"]
