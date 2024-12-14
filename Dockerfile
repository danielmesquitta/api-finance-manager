FROM golang:1.23.4-alpine as builder
WORKDIR /app
COPY . .
# Install migration tool
RUN go install github.com/steebchen/prisma-client-go@latest
# Migrate
RUN prisma-client-go migrate deploy --schema=sql/schema.prisma
# Build
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ./tmp/restapi ./cmd/restapi

FROM scratch
COPY --from=builder /app/tmp/restapi .
EXPOSE 8080
CMD ["./restapi"]
