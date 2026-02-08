package handlers

import (
	"errors"
	"net/http"

	"music-service/services"

	"github.com/gin-gonic/gin"
)

type ArtistHandler struct {
	ArtistSvc services.ArtistService
}

func NewArtistHandler(svc services.ArtistService) *ArtistHandler {
	return &ArtistHandler{ArtistSvc: svc}
}

// GET /artist/me (auth + role=artist)
func (h *ArtistHandler) Me(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	a, err := h.ArtistSvc.Me(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "artist profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get artist profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": a})
}

type upsertArtistRequest struct {
	Name       string  `json:"name" binding:"required"`
	Bio        string  `json:"bio"`
	AvatarPath *string `json:"avatar_path"` // берёшь из /uploads/covers
}

// POST /artist/profile (auth + role=artist)
func (h *ArtistHandler) CreateProfile(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req upsertArtistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := h.ArtistSvc.UpsertMe(c.Request.Context(), userID, req.Name, req.Bio, req.AvatarPath)
	if err != nil {
		if errors.Is(err, services.ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": "artist profile already exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": a})
}
