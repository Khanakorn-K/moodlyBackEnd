package routes

import (
	"moodly/config/initializers"
	"moodly/internal/controllers/authcontroller"
	"moodly/internal/controllers/customercontroller"
	"moodly/internal/controllers/insightcontroller"
	"moodly/internal/controllers/moodlogscontroller"
	"moodly/internal/controllers/overviewcontroller"
	"moodly/internal/middlewares"
	"moodly/internal/repositoriesImpl"
	"moodly/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	registerAuthRoutes(r)
	registerMoodLogsRoutes(r)
	registerCustomCauseRoutes(r)
	registerInsightRoutes(r)
	registerOverviewRoutes(r)
}

func registerAuthRoutes(r *gin.Engine) {
	authRepo := repositoriesImpl.NewAuthRepositoryImpl(initializers.DB)
	authService := services.NewAuthService(authRepo)
	authController := authcontroller.NewAuthController(authService)

	auth := r.Group("/auth")

	auth.POST("/register", authController.HandleRegister)
	auth.POST("/login", authController.HandleLogin)

	// สำหรับ Go OAuth redirect flow
	// auth.POST("/google/login", authController.HandleGoogleLogin)
	// auth.GET("/google/callback", authController.HandleGoogleCallback)

	// สำหรับ NextAuth ส่ง profile จาก Google มาให้ backend
	auth.POST("/oauth/google", authController.HandleOAuthGoogleLogin)
}

func registerMoodLogsRoutes(r *gin.Engine) {
	moodLogsRepo := repositoriesImpl.NewMoodLogsRepository(initializers.DB)
	moodLogsService := services.NewMoodLogsService(moodLogsRepo)
	moodLogsController := moodlogscontroller.NewMoodLogsController(moodLogsService)

	mood := r.Group("/mood-logs")
	mood.Use(middlewares.AuthMiddleware())

	mood.POST("/create-mood-log", moodLogsController.CreateMoodLog)
	mood.GET("/get-mood-logs", moodLogsController.GetMoodLogsByDate)
	mood.PATCH("/update-mood-log/:id", moodLogsController.UpdateMoodLog)
	mood.DELETE("/delete-mood-log/:id", moodLogsController.DeleteMoodLog)
}

func registerCustomCauseRoutes(r *gin.Engine) {
	customCauseRepo := repositoriesImpl.NewCustomCauseRepositoryImpl(initializers.DB)
	customCauseService := services.NewCustomCauseService(customCauseRepo)
	customCauseController := customercontroller.NewCustomCauseController(customCauseService)

	cause := r.Group("/custom-causes")
	cause.Use(middlewares.AuthMiddleware())

	cause.POST("/create-custom-cause", customCauseController.CreateCause)
	cause.GET("/get-custom-causes", customCauseController.GetCauses)
	cause.PATCH("/update-custom-cause/:id", customCauseController.UpdateCause)
	cause.DELETE("/delete-custom-cause/:id", customCauseController.DeleteCause)
}

func registerInsightRoutes(r *gin.Engine) {
	insightRepo := repositoriesImpl.NewInsightRepositoryImpl(initializers.DB)
	insightService := services.NewInsightService(insightRepo)
	insightController := insightcontroller.NewInsightController(insightService)

	insight := r.Group("/insights")
	insight.Use(middlewares.AuthMiddleware())

	insight.GET("/get-insights", insightController.GetInsights)
}

func registerOverviewRoutes(r *gin.Engine) {
	overviewRepo := repositoriesImpl.NewOverviewRepositoryImpl(initializers.DB)
	overviewService := services.NewOverviewService(overviewRepo)
	overviewController := overviewcontroller.NewOverviewController(overviewService)

	overview := r.Group("/overview")
	overview.Use(middlewares.AuthMiddleware())

	overview.GET("/get-monthly-average-mood", overviewController.GetMonthlyAverageMood)
	overview.GET("/get-overview", overviewController.GetOverview)
}
