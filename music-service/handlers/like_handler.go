package handlers

import (
	"net/http"
	"strconv"

	"music-service/services"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	LikeSvc services.LikeService
}

func NewLikeHandler(svc services.LikeService) *LikeHandler {
	return &LikeHandler{LikeSvc: svc}
}

// POST /tracks/:id/like (auth)
func (h *LikeHandler) Like(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	trackID, err := strconv.Atoi(c.Param("id"))
	if err != nil || trackID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track id"})
		return
	}

	if err := h.LikeSvc.Like(c.Request.Context(), userID, trackID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DELETE /tracks/:id/like (auth)
func (h *LikeHandler) Unlike(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	trackID, err := strconv.Atoi(c.Param("id"))
	if err != nil || trackID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track id"})
		return
	}

	if err := h.LikeSvc.Unlike(c.Request.Context(), userID, trackID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GET /users/me/likes?limit=20&offset=0 (auth)
func (h *LikeHandler) MyLikes(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit := 20
	offset := 0
	if raw := c.Query("limit"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			limit = v
		}
	}
	if raw := c.Query("offset"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			offset = v
		}
	}

	ids, err := h.LikeSvc.MyLikes(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list likes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ids})
}
