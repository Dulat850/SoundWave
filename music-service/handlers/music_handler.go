package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"music-service/services"

	"github.com/gin-gonic/gin"
)

type MusicHandler struct {
	MusicSvc services.MusicService
}

func NewMusicHandler(musicSvc services.MusicService) *MusicHandler {
	return &MusicHandler{MusicSvc: musicSvc}
}

// GET /artists?limit=20&offset=0
func (h *MusicHandler) ListArtists(c *gin.Context) {
	limit := parseIntQuery(c, "limit", 20)
	offset := parseIntQuery(c, "offset", 0)

	artists, err := h.MusicSvc.ListArtists(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list artists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": artists})
}

// GET /artists/:id
func (h *MusicHandler) GetArtistByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	artist, err := h.MusicSvc.GetArtistByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get artist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": artist})
}

func (h *MusicHandler) SearchTracks(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q is required"})
		return
	}

	tracks, err := h.MusicSvc.SearchTracks(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search tracks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tracks})
}

func parseIntQuery(c *gin.Context, key string, def int) int {
	raw := c.Query(key)
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	return v
}
