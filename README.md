### For Live Reloading: ###
1. install [Air|https://github.com/cosmtrek/air]: `go install github.com/cosmtrek/air@latest` 
2. run `air`
3. if you get `zsh: command not found: air` you need to add to PATH: `PATH=$PATH:$(go env GOPATH)/bin`

### For migrations: ### 
1. install [Goose|https://pressly.github.io/] `brew install goose`
2. Note, this binary is separate from the go package.
3. `goose postgres "host=localhost dbname=string_db user=string_db password=string_password sslmode=disable" down-to 0`

### Postgres & Redis - Docker Compose: ***local dev only*** ###
1. To build and start the docker containers for the first time: `docker-compose -f docker-compose.yml up`
2. To shutdown the docker containers press `ctl + c`
3. To start them again: `docker start string-api_db_1 string-api_redis_1`
4. To stop them: `docker stop string-api_db_1 string-api_redis_1`

### For local development: ### 
`cd api`
run `go install` to get dependencies installed

### For local testing: ###
run `go test`

### Unit21: ### 
This is a 3rd party service that offers the ability to evaluate risk at a transaction level and identify fraud. A client file exists to connect to their API. Documentation is here: https://docs.unit21.ai/reference/entities-api
You can create a test API key on the Unit21 dashboard. You will need to be setup as an Admin. Here are the instructions: https://docs.unit21.ai/reference/generate-api-keys
When setting up the production env variables, their URL will be: https://api.unit21.com/v1
