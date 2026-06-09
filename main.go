package main

import (
	"moodly/controllers"
	"moodly/initializers"
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

	// PostRepo := repositories.NewPostRepository(initializers.DB)
	// PostService := services.NewPostService(PostRepo)
	// PostController := controllers.NewPostController(PostService)

	// r.POST("/posts", PostController.CreatePost)
	// r.GET("/posts", PostController.GetPosts)
	// r.GET("/posts/:id", PostController.GetPostByID)
	// r.PUT("/posts/:id", PostController.UpdatePost)
	// r.DELETE("/posts/:id", PostController.DeletePost)

	AuthRepo := repositories.NewAuthRepository(initializers.DB)
	AuthService := services.NewAuthService(AuthRepo)
	AuthController := controllers.NewAuthController(AuthService)

	r.POST("/register", AuthController.HandleRegister)
	r.POST("/login", AuthController.HandleLogin)

	r.Run()
}
