package handlers

import "github.com/gin-gonic/gin"

func respondData(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"ok":   true,
		"data": data,
	})
}

func respondMessage(c *gin.Context, status int, message string, data any) {
	payload := gin.H{
		"ok":      true,
		"message": message,
	}
	if data != nil {
		payload["data"] = data
	}
	c.JSON(status, payload)
}

func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"ok":    false,
		"error": message,
	})
}
