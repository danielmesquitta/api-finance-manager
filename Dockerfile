FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN go install github.com/danielmesquitta/prisma-go-tools@latest
RUN go install github.com/steebchen/prisma-client-go@latest
RUN make migrate
RUN make build

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/tmp/restapi .
EXPOSE 8080
CMD ["./restapi"]
