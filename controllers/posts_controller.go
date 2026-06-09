package controllers

import (
	models "moodly/Models"
	"moodly/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	service *services.PostService
}

func NewPostController(service *services.PostService) *PostController {
	return &PostController{
		service: service,
	}
}

func (pc *PostController) CreatePost(c *gin.Context) {
	var post models.PostModel

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
			"error":   err.Error(),
		})
		return
	}

	if err := pc.service.CreatePost(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to create post",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "created",
		"data":    post,
	})
}

func (pc *PostController) GetPosts(c *gin.Context) {
	posts, err := pc.service.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get posts",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    posts,
	})
}

func (pc *PostController) GetPostByID(c *gin.Context) {
	id := c.Param("id")

	post, err := pc.service.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "post not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    post,
	})
}

func (pc *PostController) UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var body models.PostModel

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
			"error":   err.Error(),
		})
		return
	}

	post, err := pc.service.UpdatePost(id, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update post",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "updated",
		"data":    post,
	})
}

func (pc *PostController) DeletePost(c *gin.Context) {
	id := c.Param("id")

	if err := pc.service.DeletePost(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete post",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
	})
}
