package main

import (
	models "moodly/Models"
	"moodly/initializers"
)

// go run migrate/migrate.go ต้องรับคำสั่งนี้

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.PostModel{})
}
