package handlers

import (
	"net/http"

	"music-service/services"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	Storage services.StorageService
}

func NewUploadHandler(storage services.StorageService) *UploadHandler {
	return &UploadHandler{Storage: storage}
}

// POST /uploads/audio  (artist)
func (h *UploadHandler) UploadAudio(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	rel, err := h.Storage.SaveAudio(c.Request.Context(), fh)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// URL, который фронт сможет использовать, если ты раздаёшь /static
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"path": rel,
			"url":  "/static/" + rel,
		},
	})
}

// POST /uploads/covers (artist)
func (h *UploadHandler) UploadCover(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	rel, err := h.Storage.SaveCover(c.Request.Context(), fh)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"path": rel,
			"url":  "/static/" + rel,
		},
	})
}
