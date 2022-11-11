#!/bin/sh

# export env variables from .env file
export $(grep -v '^#' .env | xargs)

# run db migrations
echo "Running DB Migrations..."
DB_CONFIG="host=$DB_HOST user=$DB_USERNAME dbname=$DB_NAME sslmode=disable password=$DB_PASSWORD"
# goose postgres "$DB_CONFIG" reset # uncomment to reset db on build
cd migrations
goose postgres "$DB_CONFIG" up
echo '>> Done!'
# TODO: If migration fails then exit

cd ..
air