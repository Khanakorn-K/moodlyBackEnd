package main

import (
	"moodly/controllers"
	"moodly/initializers"

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

	r.POST("/posts", controllers.CreatePost)
	r.GET("/posts", controllers.GetPosts)
	r.GET("/posts/:id", controllers.GetPostByID)
	r.PUT("/posts/:id", controllers.UpdatePost)
	r.DELETE("/posts/:id", controllers.DeletePost)
	r.Run()
}
