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
	// moodLogs
	mood := r.Group("/mood-logs")
	mood.Use(middlewares.AuthMiddleware())
	mood.POST("/create-mood-log", MoodLogsController.CreateMoodLog)
	mood.GET("/get-mood-logs", MoodLogsController.GetMoodLogsByDate)
	mood.PATCH("/update-mood-log/:id", MoodLogsController.UpdateMoodLog)
	mood.DELETE("/delete-mood-log/:id", MoodLogsController.DeleteMoodLog)

	//customCauses
	CustomCauseRepo := repositories.NewCustomCauseRepository(initializers.DB)
	CustomCauseService := services.NewCustomCauseService(CustomCauseRepo)
	CustomCauseController := controllers.NewCustomCauseController(CustomCauseService)
	cause := r.Group("/custom-causes")
	cause.Use(middlewares.AuthMiddleware())
	cause.POST("/create-custom-cause", CustomCauseController.CreateCause)
	cause.GET("/get-custom-causes", CustomCauseController.GetCauses)
	cause.PATCH("/update-custom-cause/:id", CustomCauseController.UpdateCause)
	cause.DELETE("/delete-custom-cause/:id", CustomCauseController.DeleteCause)

	//insight
	InsightRepo := repositories.NewInsightRepository(initializers.DB)
	InsightService := services.NewInsightService(InsightRepo)
	InsightController := controllers.NewInsightController(InsightService)
	insight := r.Group("/insights")
	insight.Use(middlewares.AuthMiddleware())
	insight.GET("/get-insights", InsightController.FindMoodLogs)
	//overView
	OverviewRepo := repositories.NewOverviewRepository(initializers.DB)
	OverviewService := services.NewOverviewService(OverviewRepo)
	OverviewController := controllers.NewOverviewController(OverviewService)
	overview := r.Group("/overview")
	overview.Use(middlewares.AuthMiddleware())
	overview.GET("/get-monthly-average-mood", OverviewController.GetMonthlyAverageMood)
	overview.GET("/get-overview", OverviewController.GetOverview)
	r.Run()
}
