### For Live Reloading: ###
1. install [Air|https://github.com/cosmtrek/air]: `go install github.com/cosmtrek/air@latest` 
2. run `air`
3. if you get `zsh: command not found: air` you need to add to PATH: `PATH=$PATH:$(go env GOPATH)/bin`

### For migrations: ### 
1. install [Goose|https://pressly.github.io/] `brew install goose`
2. Note, this binary is separate from the go package.
3. `goose postgres "host=localhost dbname=string_db user=string_db password=string_password sslmode=disable" down-to 0`

### Postgres Docker Compose: ***local dev only*** ### 
1. Run `docker-compose -f docker-compose.yml up`
2. Exit press `ctl + c`

### For local development: ### 
run `go install` to get dependencies installed

### For local testing: ###
run `go test`

### Unit21: ### 
This is a 3rd party service that offers the ability to evaluate risk at a transaction level and identify fraud. A client file exists to connect to their API. Documentation is here: https://docs.unit21.ai/reference/entities-api
When setting up the production env variables, their URL will be: https://api.unit21.com/v1
