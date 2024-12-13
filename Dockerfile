FROM golang:1.23.4-alpine as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ./tmp/restapi ./cmd/restapi

FROM scratch
COPY --from=builder /app/tmp/restapi .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./restapi"]
