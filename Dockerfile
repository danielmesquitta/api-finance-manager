FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN go install github.com/steebchen/prisma-client-go@latest
RUN make migrate
RUN make build

FROM alpine:latest
COPY --from=builder /app/tmp/restapi .
RUN apk upgrade --no-cache
RUN apk add --no-cache ca-certificates
EXPOSE 8080
CMD ["./restapi"]
