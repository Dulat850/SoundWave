package handlers

import (
	"errors"
	"music-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenreHandler struct {
	GenreSvc services.GenreService
}

func NewGenreHandler(svc services.GenreService) *GenreHandler {
	return &GenreHandler{GenreSvc: svc}
}

// GET /genres?limit=50&offset=0 (public)
func (h *GenreHandler) List(c *gin.Context) {
	limit := parseInt(c.Query("limit"), 50)
	offset := parseInt(c.Query("offset"), 0)

	items, err := h.GenreSvc.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list genres"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

type createGenreRequest struct {
	Name string `json:"name" binding:"required"`
}

// POST /genres (artist)
func (h *GenreHandler) Create(c *gin.Context) {
	var req createGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g, err := h.GenreSvc.Create(c.Request.Context(), req.Name)
	if err != nil {
		if errors.Is(err, services.ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": "genre already exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": g})
}
