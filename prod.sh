#!/bin/bash

# echo "Migrating database..."
# migrate -path ./db/migrations -database "$DATABASE_URL" -verbose force 1 || true
# migrate -path ./db/migrations -database "$DATABASE_URL" -verbose up

echo "Starting DiceRealms in Production mode using Render..."
./dicerealms