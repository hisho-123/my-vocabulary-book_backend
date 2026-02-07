package controller

import (
	"backend/src/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	requestHeader := c.GetHeader("Token")
	if requestHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token is required.",
		})
		return
	}

	if err := usecase.ValidateUserToken(requestHeader); err != nil {
		c.JSON(statusCode(err), gin.H{
			"error": "Invalid or expired token.",
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
