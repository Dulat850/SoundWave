package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no current user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": currentUser})
}
