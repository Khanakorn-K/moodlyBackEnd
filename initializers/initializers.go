package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

// จัดการ env
func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}
}
