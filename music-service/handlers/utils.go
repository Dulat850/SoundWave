package handlers

import (
	"music-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseInt(raw string, def int) int {
	if raw == "" {
		return def
	}

	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}

	return v
}

func currentUserID(c *gin.Context) (int, bool) {
	v, ok := c.Get("currentUser")
	if !ok || v == nil {
		return 0, false
	}

	u, ok := v.(*models.User)
	if !ok || u == nil || u.ID <= 0 {
		return 0, false
	}

	return u.ID, true
}
