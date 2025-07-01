FROM golang:1.24 AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /driver-location-api ./cmd/main.go

FROM scratch

COPY --from=builder /driver-location-api /driver-location-api

ENV APP_PORT=<your port> \
    MONGO_URI=<your mongo uri> \
    MONGO_DB=<your mongo database name> \
    INTERNAL_API_KEY=<your api key> \
    LOG_LEVEL=info

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD \
  wget -qO- http://localhost:${APP_PORT}/v1/healthz || exit 1

ENTRYPOINT ["/driver-location-api"]
