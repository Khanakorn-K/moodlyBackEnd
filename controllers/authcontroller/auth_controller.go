package authcontroller

import (
	"crypto/subtle"
	models "moodly/Models"
	"moodly/services"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

const googleOAuthStateCookie = "google_oauth_state"

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(services *services.AuthService) *AuthController {
	return &AuthController{service: services}
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

func (ac *AuthController) HandleGoogleLogin(c *gin.Context) {
	state, err := ac.service.GenerateOAuthState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to start google login",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(googleOAuthStateCookie, state, 600, "/", "", false, true)

	url := ac.service.GetGoogleOAuthURL(state)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (ac *AuthController) HandleGoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "code is required",
		})
		return
	}

	stateCookie, err := c.Cookie(googleOAuthStateCookie)
	if err != nil || state == "" || subtle.ConstantTimeCompare([]byte(state), []byte(stateCookie)) != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid oauth state",
		})
		return
	}
	c.SetCookie(googleOAuthStateCookie, "", -1, "/", "", false, true)

	token, err := ac.service.LoginWithGoogle(code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	redirectURL := os.Getenv("FRONTEND_AUTH_CALLBACK_URL")
	if redirectURL == "" {
		redirectURL = "http://localhost:3000/auth/callback"
	}

	parsedRedirectURL, err := url.Parse(redirectURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "invalid frontend callback url",
		})
		return
	}

	query := parsedRedirectURL.Query()
	query.Set("token", token)
	parsedRedirectURL.RawQuery = query.Encode()

	c.Redirect(
		http.StatusTemporaryRedirect,
		parsedRedirectURL.String(),
	)
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
