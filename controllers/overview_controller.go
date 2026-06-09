package controllers

import (
	"moodly/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OverviewController struct {
	service *services.OverviewService
}

func NewOverviewController(service *services.OverviewService) *OverviewController {
	return &OverviewController{service: service}
}

func (oc *OverviewController) GetMonthlyAverageMood(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	month := c.Query("month")

	result, err := oc.service.GetMonthlyAverageMood(userID, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (oc *OverviewController) GetOverview(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	result, err := oc.service.GetOverview(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
