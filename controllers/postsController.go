package controllers

import (
	models "moodly/Models"
	"moodly/initializers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	var post models.PostModel

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
			"error":   err.Error(),
		})
		return
	}

	result := initializers.DB.Create(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create post",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "created",
		"data":    post,
	})
}

func GetPosts(c *gin.Context) {
	var posts []models.PostModel

	result := initializers.DB.Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get posts",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    posts,
	})
}

func GetPostByID(c *gin.Context) {
	id := c.Param("id")

	var post models.PostModel

	result := initializers.DB.First(&post, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "post not found",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    post,
	})
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var post models.PostModel
	//หา ก่อน
	result := initializers.DB.First(&post, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "post not found",
			"error":   result.Error.Error(),
		})
		return
	}

	var body models.PostModel
	// เช้ค req
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
			"error":   err.Error(),
		})
		return
	}

	result = initializers.DB.Model(&post).Updates(body)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update post",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "updated",
		"data":    post,
	})
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")

	var post models.PostModel

	result := initializers.DB.First(&post, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "post not found",
			"error":   result.Error.Error(),
		})
		return
	}

	result = initializers.DB.Delete(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete post",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
	})
}
