FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o build order-service/cmd/order
RUN go build -o migrator order-service/cmd/migrator

FROM golang:1.21
COPY --from=builder /app/build /build
COPY --from=builder /app/migrator /migrator
CMD ["/build"]