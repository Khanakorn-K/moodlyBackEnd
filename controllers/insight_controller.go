package controllers

import (
	"errors"
	"moodly/helpers"
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

func (ic *InsightController) GetInsights(c *gin.Context) {
	userID, ok := helpers.GetUserIDFromContext(c)
	if !ok {
		helpers.CreateAPIErrorResponse(
			c,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"unauthorized",
		)
		return
	}

	selectedDate := c.Query("selectedDate")

	result, err := ic.service.GetInsights(userID, selectedDate)
	if err != nil {
		if errors.Is(err, services.ErrInvalidDateFormat) {
			helpers.CreateAPIErrorResponse(
				c,
				http.StatusBadRequest,
				"INVALID_DATE_FORMAT",
				err.Error(),
			)
			return
		}

		helpers.CreateAPIErrorResponse(
			c,
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"internal server error",
		)
		return
	}

	helpers.CreateAPIResponse(c, http.StatusOK, result)
}
