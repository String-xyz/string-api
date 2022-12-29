package scripts

import (
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/joho/godotenv"
)

func GenerateWallet() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	err = service.GenerateWallet()
	if err != nil {
		panic(err)
	}
}
