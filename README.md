For Live Reloading: 
1. install [Air|https://github.com/cosmtrek/air]: `go install github.com/cosmtrek/air@latest` 
2. run `air`

For migrations:
1. install [Goose|https://pressly.github.io/] `brew install goose`
2. Note, this binary is separate from the go package.

Postgres Docker Compose: ***local dev only***
1. Run `docker-compose -f docker-compose.yml up`
2. Exit press `ctl + c`