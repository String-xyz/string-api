package main

import (
	"os"

	"github.com/String-xyz/string-api/scripts"
)

func main() {
	var script string
	if len(os.Args) > 1 {
		script = os.Args[1]
	}

	if script == "generate_wallet" {
		scripts.GenerateWallet()
	} // etc...
}
