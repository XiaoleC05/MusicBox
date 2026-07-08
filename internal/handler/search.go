package handler

import (
	"net/http"
	"strconv"

	"github.com/XiaoleC05/MusicBox/internal/adapter"
	"github.com/gin-gonic/gin"
)

var kugouAdapter *adapter.KugouAdapter

func InitAdapters(kugouCookie string) {
	kugouAdapter = adapter.NewKugouAdapter(kugouCookie)
}

func Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	platform := c.Query("platform")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var songs []adapter.Song
	var err error

	if platform == "" {
		if kugouAdapter != nil && kugouAdapter.IsAvailable() {
			songs, err = kugouAdapter.Search(query, page)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	} else if platform == "kugou" {
		if kugouAdapter == nil || !kugouAdapter.IsAvailable() {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "kugou adapter not available"})
			return
		}
		songs, err = kugouAdapter.Search(query, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported platform"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

func Play(c *gin.Context) {
	platform := c.Param("platform")
	songID := c.Param("songId")
	quality := c.DefaultQuery("quality", "standard")

	var playURL string
	var err error

	if platform == "kugou" {
		if kugouAdapter == nil || !kugouAdapter.IsAvailable() {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "kugou adapter not available"})
			return
		}
		playURL, err = kugouAdapter.GetPlayURL(songID, quality)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported platform"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"play_url": playURL,
		"platform": platform,
		"song_id":  songID,
		"quality":  quality,
	})
}
