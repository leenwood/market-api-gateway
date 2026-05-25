FROM golang:1.26 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/gateway ./cmd/gateway

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/bin/gateway /gateway

EXPOSE 80

ENTRYPOINT ["/gateway"]
