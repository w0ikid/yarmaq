FROM golang:1.25.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make curl
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy shared pkg first to leverage cache
COPY pkg ./pkg

# Copy account service
COPY apps/accounts-service ./apps/accounts-service

WORKDIR /app/apps/accounts-service
RUN go mod download
RUN go build -o main ./cmd/api/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata postgresql-client curl

COPY --from=builder /app/apps/accounts-service/main .
COPY --from=builder /app/apps/accounts-service/migrations ./migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Use shared start script from deployment
COPY deployment/scripts/start.sh /app/start.sh
RUN chmod +x /app/start.sh

# secrets (if needed, but better via volumes or env)
# COPY secrets/zitadel.json /app/secrets/zitadel.json

CMD ["/app/start.sh"]
