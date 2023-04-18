#!/bin/sh

# export env variables from .env file
export $(grep -v '^#' .env | xargs)

if [ "$DB_RESET" = "true" ]; then
  # run db migrations
  echo "----- Running migrations..."
  cd migrations

  DB_CONFIG="host=$DB_HOST user=$DB_USERNAME dbname=$DB_NAME sslmode=disable password=$DB_PASSWORD"
  goose postgres "$DB_CONFIG" reset
  goose postgres "$DB_CONFIG" up
  cd ..
  echo "----- ...Migrations done"
  echo "----- Seeding data..."
  go run script.go data_seeding local
  echo "----- ...Data seeded"
fi

# run app
if [ "$DEBUG_MODE" = "true" ]; then
  echo "----- DEBUG_MODE is true"
  air -c .air-debug.toml
else
  echo "----- DEBUG_MODE is false"
  air
fi