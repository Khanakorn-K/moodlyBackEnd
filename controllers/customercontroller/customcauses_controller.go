package customercontroller

import (
	models "moodly/Models"
	"moodly/helpers"
	"moodly/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomCauseController struct {
	service *services.CustomCauseService
}

func NewCustomCauseController(service *services.CustomCauseService) *CustomCauseController {
	return &CustomCauseController{service: service}
}

func (cc *CustomCauseController) CreateCause(c *gin.Context) {
	userID, ok := helpers.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	var req createCustomCausesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	customcauses := models.CustomCause{
		UserID: userID,
		Name:   req.Name,
	}

	if err := cc.service.CreateCause(&customcauses); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "cause created",
		"cause":   customcauses.Name,
	})
}

func (cc *CustomCauseController) GetCauses(c *gin.Context) {
	userID, ok := helpers.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	causes, err := cc.service.GetCauses(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"causes": causes,
	})
}

func (cc *CustomCauseController) UpdateCause(c *gin.Context) {
	userID, ok := helpers.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid cause id"})
		return
	}

	var req createCustomCausesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	cause := models.CustomCause{
		ID:     uint(id),
		UserID: userID,
		Name:   req.Name,
	}

	if err := cc.service.UpdateCause(&cause); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "cause updated",
		"cause":   cause,
	})
}

func (cc *CustomCauseController) DeleteCause(c *gin.Context) {
	userID, ok := helpers.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid cause id"})
		return
	}

	if err := cc.service.DeleteCause(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "cause deleted",
	})
}
