package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"music-service/services"

	"github.com/gin-gonic/gin"
)

type PlaylistHandler struct {
	PlaylistSvc services.PlaylistService
}

func NewPlaylistHandler(svc services.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{PlaylistSvc: svc}
}

type createPlaylistRequest struct {
	Name string `json:"name" binding:"required"`
}

// POST /playlists (auth)
func (h *PlaylistHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := h.PlaylistSvc.Create(c.Request.Context(), userID, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": p})
}

// GET /users/me/playlists (auth)
func (h *PlaylistHandler) MyPlaylists(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	items, err := h.PlaylistSvc.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list playlists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GET /playlists/:id
func (h *PlaylistHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	p, err := h.PlaylistSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get playlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": p})
}

type renamePlaylistRequest struct {
	Name string `json:"name" binding:"required"`
}

// PATCH /playlists/:id (auth)
func (h *PlaylistHandler) Rename(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req renamePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.PlaylistSvc.Rename(c.Request.Context(), userID, id, req.Name)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rename playlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type addTrackRequest struct {
	TrackID int `json:"track_id" binding:"required"`
}

// POST /playlists/:id/tracks (auth)
func (h *PlaylistHandler) AddTrack(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	playlistID, err := strconv.Atoi(c.Param("id"))
	if err != nil || playlistID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req addTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.TrackID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "track_id is required"})
		return
	}

	err = h.PlaylistSvc.AddTrack(c.Request.Context(), userID, playlistID, req.TrackID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add track"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DELETE /playlists/:id/tracks/:trackId (auth)
func (h *PlaylistHandler) RemoveTrack(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	playlistID, err := strconv.Atoi(c.Param("id"))
	if err != nil || playlistID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	trackID, err := strconv.Atoi(c.Param("trackId"))
	if err != nil || trackID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trackId"})
		return
	}

	err = h.PlaylistSvc.RemoveTrack(c.Request.Context(), userID, playlistID, trackID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove track"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
