#!/bin/bash
set -e

if [ -z "$DATABASE_URL" ]; then
    DATABASE_URL="postgresql://postgres:postgres@localhost:5480/alizia?sslmode=disable"
fi

echo "Running migrations..."
migrate -path db/migrations -database "$DATABASE_URL" up
echo "Migrations complete."
