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
	//ถ้ามีการเพิ่ม model หรือแก้ไข อย่าลืม migrate
	initializers.DB.AutoMigrate(
		&models.User{},
		&models.OAuthAccount{},
		&models.MoodLog{},
		&models.CustomCause{},
	)
}
