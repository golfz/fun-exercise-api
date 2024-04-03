FROM golang:latest as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v ./...

RUN go build -o ./out/wallet-api .

# ========================

FROM alpine:latest

COPY --from=build-base /app/out/wallet-api /app/wallet-api

CMD ["/app/wallet-api"]

