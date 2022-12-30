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

	if script == "data_seeding" {
		dataSeedingArgs := "local"
		if len(os.Args) > 2 {
			dataSeedingArgs = os.Args[2]
		}
		if dataSeedingArgs == "local" {
			scripts.MockSeeding()
		} else {
			scripts.DataSeeding()
		}
	} else if script == "generate_wallet" {
		scripts.GenerateWallet()
	} // etc...
}
