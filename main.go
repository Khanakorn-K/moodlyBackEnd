package main

import (
	"moodly/controllers"
	"moodly/initializers"
	"moodly/middlewares"
	"moodly/repositories"
	"moodly/services"

	"github.com/gin-gonic/gin"
)

// init() เป็น function พิเศษของ Go
// มันจะถูกรัน อัตโนมัติก่อน main()
func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {

	r := gin.Default()

	//auth
	AuthRepo := repositories.NewAuthRepository(initializers.DB)
	AuthService := services.NewAuthService(AuthRepo)
	AuthController := controllers.NewAuthController(AuthService)
	auth := r.Group("/auth")
	auth.POST("/register", AuthController.HandleRegister)
	auth.POST("/login", AuthController.HandleLogin)
	//moodLogs
	MoodLogsRepo := repositories.NewMoodLogsRepository(initializers.DB)
	MoodLogsService := services.NewMoodLogsService(MoodLogsRepo)
	MoodLogsController := controllers.NewMoodLogsController(MoodLogsService)
	mood := r.Group("/mood")
	mood.Use(middlewares.AuthMiddleware())
	mood.POST("/createmoodlog", MoodLogsController.CreateMoodLog)
	mood.GET("/getmoodlogs", MoodLogsController.GetMoodLogsByDate)
	mood.PATCH("/updatemoodlog/:id", MoodLogsController.UpdateMoodLog)
	mood.DELETE("/deletemoodlog/:id", MoodLogsController.DeleteMoodLog)
	r.Run()
}
