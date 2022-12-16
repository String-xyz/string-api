#!/bin/sh

# export env variables from .env file
export $(grep -v '^#' .env | xargs)

# run db migrations
echo "----- Running migrations..."
cd migrations
cd mocks

DB_CONFIG="host=$DB_HOST user=$DB_USERNAME dbname=$DB_NAME sslmode=disable password=$DB_PASSWORD"
goose postgres "$DB_CONFIG" reset
cd ..
goose postgres "$DB_CONFIG" reset
goose postgres "$DB_CONFIG" up
cd ..
echo "----- ...Migrations done"
echo "----- Seeding data..."
go run scripts/data_seeding.go local
echo "----- ...Data seeded"

# run app
air