For Live Reloading: 
1. install [Air|https://github.com/cosmtrek/air]: `go install github.com/cosmtrek/air@latest` 
2. run `air`
3. if you get `zsh: command not found: air` you need to add to PATH: `PATH=$PATH:$(go env GOPATH)/bin`

For migrations:
1. install [Goose|https://pressly.github.io/] `brew install goose`
2. Note, this binary is separate from the go package.

