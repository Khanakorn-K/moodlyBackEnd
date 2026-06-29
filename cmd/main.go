// @title Moodly API
// @version 1.0
// @description Moodly Backend API
// @host 54.251.189.159:8080 // ถ้า เปิดปิด instance อย่าลืมมาแก้
// @BasePath /

package main

import (
	_ "moodly/docs"
	"os"

	"moodly/config/initializers"
	"moodly/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
	initializers.ConnectKafka(os.Getenv("KAFKA_ADDR"), "moodlyTP")
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			os.Getenv("URL_SERVER"),
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))
	routes.RegisterRoutes(r)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
