package controller

import (
	"backend/src/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBookListHandler(c *gin.Context) {
	requestHeader := c.GetHeader("Token")
	if requestHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token is required.",
		})
		return
	}

	bookList, err := usecase.GetBookList(requestHeader)
	if err != nil {
		c.JSON(statusCode(err), gin.H{
			"error": "Could not get book list.",
		})
		return
	}

	c.JSON(http.StatusOK, bookList)
}

func GetBookHandler(c *gin.Context) {
	requestHeader := c.GetHeader("Token")
	if requestHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token is required.",
		})
		return
	}

	bookIdStr := c.Query("bookId")
	bookId, err := strconv.Atoi(bookIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book ID format.",
		})
		return
	}

	book, err := usecase.GetBook(requestHeader, bookId)
	if err != nil {
		c.JSON(statusCode(err), gin.H{
			"error": "Could not get book.",
		})
		return
	}

	c.JSON(http.StatusOK, book)
}
