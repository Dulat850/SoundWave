package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"music-service/repositories"
	"music-service/services"

	"github.com/gin-gonic/gin"
)

type AlbumHandler struct {
	AlbumSvc   services.AlbumService
	ArtistRepo repositories.ArtistRepository
}

func NewAlbumHandler(albumSvc services.AlbumService, artistRepo repositories.ArtistRepository) *AlbumHandler {
	return &AlbumHandler{AlbumSvc: albumSvc, ArtistRepo: artistRepo}
}

type createAlbumRequest struct {
	Title      string  `json:"title" binding:"required"`
	CoverPath  *string `json:"cover_path"`  // путь из /uploads/covers
	ReleasedAt *string `json:"released_at"` // YYYY-MM-DD
}

// POST /artist/albums (auth + role=artist)
func (h *AlbumHandler) CreateAsArtist(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	artist, err := h.ArtistRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artist profile not found"})
		return
	}

	var req createAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := h.AlbumSvc.Create(c.Request.Context(), artist.ID, req.Title, req.CoverPath, req.ReleasedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": a})
}

// GET /albums/:id (public)
func (h *AlbumHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	a, err := h.AlbumSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": a})
}

// GET /artists/:id/albums?limit=20&offset=0 (public)
func (h *AlbumHandler) ListByArtist(c *gin.Context) {
	artistID, err := strconv.Atoi(c.Param("id"))
	if err != nil || artistID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}

	limit := parseInt(c.Query("limit"), 20)
	offset := parseInt(c.Query("offset"), 0)

	items, err := h.AlbumSvc.ListByArtistID(c.Request.Context(), artistID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list albums"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}
