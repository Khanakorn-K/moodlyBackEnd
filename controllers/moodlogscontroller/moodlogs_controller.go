package moodlogscontroller

import (
	models "moodly/Models"
	"moodly/helpers"
	"moodly/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MoodLogsController struct {
	service *services.MoodLogsService
}

func NewMoodLogsController(service *services.MoodLogsService) *MoodLogsController {
	return &MoodLogsController{service: service}
}

func (mc *MoodLogsController) CreateMoodLog(c *gin.Context) {
	userID, ok := helpers.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	var req moodlogsCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	moodLog := models.MoodLog{
		UserID: userID,
		Mood:   req.Mood,
		Note:   req.Note,
		Causes: req.Causes,
	}

	if err := mc.service.CreateMoodLog(&moodLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "mood log created",
		"mood_log": moodLog,
	})
}

func (mc *MoodLogsController) GetMoodLogsByDate(c *gin.Context) {
	userID, ok := helpers.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	date := c.Query("date")

	moodLogs, err := mc.service.GetMoodLogsByDate(userID, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mood_logs": moodLogs,
	})
}

func (mc *MoodLogsController) UpdateMoodLog(c *gin.Context) {
	userID, ok := helpers.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid mood log id"})
		return
	}

	var req moodlogsCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	moodLog := models.MoodLog{
		ID:     uint(id),
		UserID: uint(userID),
		Mood:   req.Mood,
		Note:   req.Note,
		Causes: req.Causes,
	}

	if err := mc.service.UpdateMoodLog(&moodLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "mood log updated",
		"mood_log": &moodLog,
	})
}

func (mc *MoodLogsController) DeleteMoodLog(c *gin.Context) {
	userID, ok := helpers.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid mood log id"})
		return
	}

	if err := mc.service.DeleteMoodLog(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "mood log deleted",
	})
}
