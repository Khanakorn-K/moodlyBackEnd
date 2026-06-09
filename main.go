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

	AuthRepo := repositories.NewAuthRepository(initializers.DB)
	AuthService := services.NewAuthService(AuthRepo)
	AuthController := controllers.NewAuthController(AuthService)

	auth := r.Group("/auth")

	auth.POST("/register", AuthController.HandleRegister)
	auth.POST("/login", AuthController.HandleLogin)

	PostRepo := repositories.NewPostRepository(initializers.DB)
	PostService := services.NewPostService(PostRepo)
	PostController := controllers.NewPostController(PostService)

	post := r.Group("/post")
	post.Use(middlewares.AuthMiddleware())
	post.POST("/createpost", PostController.CreatePost)
	post.GET("/getposts", PostController.GetPosts)
	post.GET("/getpost/:id", PostController.GetPostByID)
	post.PUT("/updatepost/:id", PostController.UpdatePost)
	post.DELETE("/deletepost/:id", PostController.DeletePost)

	r.Run()
}
