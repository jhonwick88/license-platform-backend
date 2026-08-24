package middleware

import "github.com/gin-gonic/gin"

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"status": "success",
		"data":   data,
	})
}

func ErrorResponse(c *gin.Context, code int, errorCode string, message string) {
	c.JSON(code, gin.H{
		"status":     "error",
		"error_code": errorCode,
		"message":    message,
	})
}
