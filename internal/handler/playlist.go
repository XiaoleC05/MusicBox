package handler

import (
	"net/http"
	"strconv"

	"github.com/XiaoleC05/MusicBox/internal/db"
	"github.com/XiaoleC05/MusicBox/internal/model"
	"github.com/gin-gonic/gin"
)

type PlaylistHandler struct {
	repo *db.PlaylistRepository
}

func NewPlaylistHandler() *PlaylistHandler {
	return &PlaylistHandler{
		repo: db.NewPlaylistRepository(),
	}
}

func (h *PlaylistHandler) List(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	playlists, err := h.repo.ListByUser(c.Request.Context(), userID)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"playlists": playlists})
}

func (h *PlaylistHandler) Create(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playlist, err := h.repo.Create(c.Request.Context(), userID, req.Name)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"playlist": playlist})
}

func (h *PlaylistHandler) Delete(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	playlistID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), playlistID, userID); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "playlist deleted"})
}

func (h *PlaylistHandler) AddSong(c *gin.Context) {
	playlistID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist ID"})
		return
	}

	var song model.PlaylistSong
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	if err := h.repo.AddSong(c.Request.Context(), playlistID, userID, song); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "song added"})
}

func (h *PlaylistHandler) RemoveSong(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	songID, err := strconv.ParseInt(c.Param("songId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid song ID"})
		return
	}

	if err := h.repo.RemoveSong(c.Request.Context(), songID, userID); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "song removed"})
}

func (h *PlaylistHandler) ListSongs(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	playlistID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist ID"})
		return
	}

	songs, err := h.repo.ListSongs(c.Request.Context(), playlistID, userID)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"songs": songs})
}
