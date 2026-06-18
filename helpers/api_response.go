package helpers

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
	Error      *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func CreateAPIResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success:    true,
		StatusCode: statusCode,
		Data:       data,
		Error:      nil,
	})
}

func CreateAPIErrorResponse(c *gin.Context, statusCode int, code string, message string) {
	c.JSON(statusCode, APIResponse{
		Success:    false,
		StatusCode: statusCode,
		Data:       nil,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}
