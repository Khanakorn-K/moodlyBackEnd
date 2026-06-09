package controllers

import (
	"moodly/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InsightController struct {
	service *services.InsightService
}

func NewInsightController(service *services.InsightService) *InsightController {
	return &InsightController{service: service}
}

func (ic *InsightController) FindMoodLogs(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	mood := c.Query("mood")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	result, err := ic.service.FindMoodLogs(userID, mood, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
