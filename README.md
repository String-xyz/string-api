### To run the APIS: ###
1. Ensure you have the infra repo where the docker compose file is now located
2. run `docker compose -f ../infra/local/docker-compose.yml up`
3. You can also run with the -d flag to keep the process in the background, ie `docker compose -f -d ../infra/local/docker-compose.yml up`

### To get a live terminal output from any repository being run by the infra docker compose: ###
1. Run `docker logs [docker container name] -f`
2. ie `docker logs platform-admin-api -f`
3. You can get the container names using `docker ps`
4. If you don't want the output in real time, you can omit the `-f` flag

### For migrations: ### 
1. This is now handled by the docker compose file

### Docker Issues?
1. If docker is giving you an error when you try to `docker-compose up --build` try the following commands in order:
2. `docker-compose down`
3. `docker system prune` (/y)
4. `docker volume prune` (/y)
5. `docker compose up --build`

### For local development: ### 
`cd api`
run `go install` to get dependencies installed

### For local testing: ###
run `go test`
or if you want to run a specific test, use `go test -run [TestName] [./path/to/dir] -v -count 1`
ie `go test -run TestGetSwapPayload ./pkg/service -v -count 1`

### To run an individual test in verbose mode ###
`go test -run [TestName] [TestDirectory] -v -count 1`
i.e. `go test -run TestGetSwapPayload ./pkg/service -v -count 1`

### Unit21: ### 
This is a 3rd party service that offers the ability to evaluate risk at a transaction level and identify fraud. A client file exists to connect to their API. Documentation is here: https://docs.unit21.ai/reference/entities-api
You can create a test API key on the Unit21 dashboard. You will need to be setup as an Admin. Here are the instructions: https://docs.unit21.ai/reference/generate-api-keys
When setting up the production env variables, their URL will be: https://api.unit21.com/v1