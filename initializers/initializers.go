package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

// จัดการ env
func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("env error : ", err)
	}
}
