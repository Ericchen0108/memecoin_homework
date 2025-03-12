package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidateID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return "", false
	}

	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return "", false
	}

	return id, true
}
