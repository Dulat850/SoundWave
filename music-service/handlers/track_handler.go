package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"music-service/models"
	"music-service/repositories"
	"music-service/services"

	"github.com/gin-gonic/gin"
)

type TrackHandler struct {
	TrackSvc   services.TrackService
	ArtistRepo repositories.ArtistRepository
}

func NewTrackHandler(trackSvc services.TrackService, artistRepo repositories.ArtistRepository) *TrackHandler {
	return &TrackHandler{TrackSvc: trackSvc, ArtistRepo: artistRepo}
}

type createTrackRequest struct {
	Title           string  `json:"title" binding:"required"`
	DurationSeconds int     `json:"duration_seconds"`
	AlbumID         *int    `json:"album_id"`
	GenreID         *int    `json:"genre_id"`
	AudioPath       string  `json:"audio_path" binding:"required"` // путь из UploadAudio
	CoverPath       *string `json:"cover_path"`                    // путь из UploadCover (опционально)
}

// POST /artist/tracks (auth + role=artist)
func (h *TrackHandler) CreateAsArtist(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	artist, err := h.ArtistRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artist profile not found"})
		return
	}

	t := &models.Track{
		ArtistID:        artist.ID,
		AlbumID:         req.AlbumID,
		GenreID:         req.GenreID,
		Title:           req.Title,
		DurationSeconds: req.DurationSeconds,
		AudioPath:       req.AudioPath,
		CoverPath:       req.CoverPath,
	}

	created, err := h.TrackSvc.Create(c.Request.Context(), t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": created})
}

// GET /tracks/:id
func (h *TrackHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	t, err := h.TrackSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get track"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": t})
}

// GET /tracks/top?limit=20
func (h *TrackHandler) Top(c *gin.Context) {
	limit := 20
	if raw := c.Query("limit"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			limit = v
		}
	}

	items, err := h.TrackSvc.Top(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get top tracks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *TrackHandler) List(c *gin.Context) {
	limit := parseInt(c.Query("limit"), 20)
	offset := parseInt(c.Query("offset"), 0)
	sort := c.DefaultQuery("sort", "new")

	items, err := h.TrackSvc.List(c.Request.Context(), limit, offset, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tracks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GET /artists/:id/tracks?limit=20&offset=0&sort=new|popular|old
func (h *TrackHandler) ListByArtist(c *gin.Context) {
	artistID, err := strconv.Atoi(c.Param("id"))
	if err != nil || artistID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}

	limit := parseInt(c.Query("limit"), 20)
	offset := parseInt(c.Query("offset"), 0)
	sort := c.DefaultQuery("sort", "new")

	items, err := h.TrackSvc.ListByArtistID(c.Request.Context(), artistID, limit, offset, sort)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GET /albums/:id/tracks?limit=20&offset=0&sort=new|popular|old
func (h *TrackHandler) ListByAlbum(c *gin.Context) {
	albumID, err := strconv.Atoi(c.Param("id"))
	if err != nil || albumID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid album id"})
		return
	}

	limit := parseInt(c.Query("limit"), 20)
	offset := parseInt(c.Query("offset"), 0)
	sort := c.DefaultQuery("sort", "new")

	items, err := h.TrackSvc.ListByAlbumID(c.Request.Context(), albumID, limit, offset, sort)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

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
