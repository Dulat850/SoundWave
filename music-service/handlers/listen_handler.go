package handlers

import (
	"net/http"
	"strconv"

	"music-service/services"

	"github.com/gin-gonic/gin"
)

type ListenHandler struct {
	ListenSvc services.ListenService
}

func NewListenHandler(svc services.ListenService) *ListenHandler {
	return &ListenHandler{ListenSvc: svc}
}

// POST /tracks/:id/listen (auth)
func (h *ListenHandler) Record(c *gin.Context) {
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

	err = h.ListenSvc.RecordListen(c.Request.Context(), &userID, trackID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GET /users/me/history?limit=20&offset=0 (auth)
func (h *ListenHandler) MyHistory(c *gin.Context) {
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

	items, err := h.ListenSvc.MyHistory(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}
