#!/bin/bash
set -e

if [ -z "$DATABASE_URL" ]; then
    DATABASE_URL="postgresql://postgres:postgres@localhost:5480/alizia?sslmode=disable"
fi

echo "Seeding database..."
psql "$DATABASE_URL" -f db/seeds/seed.sql
echo "Seeding complete."
