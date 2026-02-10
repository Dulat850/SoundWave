package handlers

import (
	"music-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	SearchSvc services.SearchService
}

func NewSearchHandler(svc services.SearchService) *SearchHandler {
	return &SearchHandler{SearchSvc: svc}
}

// GET /search?q=...&type=all|tracks|artists|genres&limit=20&offset=0
func (h *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")
	searchType := c.DefaultQuery("type", "all")

	limit := parseInt(c.Query("limit"), 20)
	offset := parseInt(c.Query("offset"), 0)

	switch searchType {
	case "all":
		tracks, artists, genres, err := h.SearchSvc.All(c.Request.Context(), q, limit, offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": gin.H{
			"tracks":  tracks,
			"artists": artists,
			"genres":  genres,
		}})
		return

	case "tracks":
		tracks, err := h.SearchSvc.Tracks(c.Request.Context(), q, limit, offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": tracks})
		return

	case "artists":
		artists, err := h.SearchSvc.Artists(c.Request.Context(), q, limit, offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": artists})
		return

	case "genres":
		genres, err := h.SearchSvc.Genres(c.Request.Context(), q, limit, offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": genres})
		return

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type"})
		return
	}
}
