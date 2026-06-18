package authcontroller

import (
	models "moodly/Models"
	"moodly/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (ac *AuthController) HandleRegister(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: &req.Password,
	}

	if err := ac.service.CreateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func (ac *AuthController) HandleLogin(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, err := ac.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"token":   token,
	})
}

type OAuthGoogleLoginRequest struct {
	Email             string `json:"email" binding:"required"`
	Name              string `json:"name"`
	Provider          string `json:"provider"`
	ProviderAccountID string `json:"providerAccountId" binding:"required"`
}

func (ac *AuthController) HandleOAuthGoogleLogin(c *gin.Context) {
	var req OAuthGoogleLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"statusCode": http.StatusBadRequest,
			"data":       nil,
			"error": gin.H{
				"code":    "INVALID_REQUEST_BODY",
				"message": err.Error(),
			},
		})
		return
	}

	token, user, err := ac.service.LoginWithOAuthGoogle(
		req.Email,
		req.Name,
		req.Provider,
		req.ProviderAccountID,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"statusCode": http.StatusUnauthorized,
			"data":       nil,
			"error": gin.H{
				"code":    "GOOGLE_LOGIN_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"statusCode": http.StatusOK,
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		},
	})
}
