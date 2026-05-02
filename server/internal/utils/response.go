package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func SuccessResponse(c *gin.Context, code int, message string, data any) {
	c.JSON(code, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string, errs any) {
	c.JSON(code, ApiResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	})
	c.Abort()
}

func ValidationErrorResponse(c *gin.Context, message string, errs any) {
	c.JSON(http.StatusBadRequest, ApiResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	})
	c.Abort()
}