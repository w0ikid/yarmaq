#!/bin/sh
set -e

DB_STRING="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_SERVICE_DB}?sslmode=disable"

echo "DB: ${DB_STRING}"
echo "Running migrations..."

GOOSE_DRIVER=postgres GOOSE_DBSTRING="${DB_STRING}" goose -dir /app/migrations up -v

echo "Starting application..."
exec /app/main