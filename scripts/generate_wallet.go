package scripts

import (
	"github.com/String-xyz/string-api/pkg/service"
	"github.com/joho/godotenv"
)

func GenerateWallet() {
	godotenv.Load(".env") // removed the err since in cloud this wont be loaded
	err := service.GenerateWallet()
	if err != nil {
		panic(err)
	}
}
