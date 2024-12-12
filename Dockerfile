FROM golang:1.23.4-alpine as builder
WORKDIR /app
COPY . .
RUN make build

FROM scratch
COPY --from=builder /app/tmp/restapi .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./restapi"]
