package handler

import (
	"net/http"
	"strconv"

	"github.com/XiaoleC05/MusicBox/internal/adapter"
	"github.com/XiaoleC05/MusicBox/internal/config"
	"github.com/XiaoleC05/MusicBox/internal/db"
	"github.com/gin-gonic/gin"
)

// getKugouAdapterForUser creates a per-request KugouAdapter from the current user's credentials.
// Returns nil if the user has no kugou credentials configured.
func getKugouAdapterForUser(c *gin.Context) (*adapter.KugouAdapter, error) {
	userID, ok := GetUserID(c)
	if !ok {
		return nil, nil
	}

	creds, err := db.NewCredentialsRepository().GetByUser(c.Request.Context(), userID)
	if err != nil {
		// No credentials record — adapter unavailable
		return nil, nil
	}

	if creds.KugouCookie == "" {
		return nil, nil
	}

	key := config.Cfg.EncryptionKey
	decrypted, err := adapter.Decrypt(creds.KugouCookie, key)
	if err != nil {
		return nil, err
	}

	return adapter.NewKugouAdapter(decrypted), nil
}

func Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	platform := c.Query("platform")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	kugouAdapter, err := getKugouAdapterForUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load credentials"})
		return
	}

	var songs []adapter.Song

	if platform == "" {
		if kugouAdapter != nil && kugouAdapter.IsAvailable() {
			songs, err = kugouAdapter.Search(query, page)
			if err != nil {
				respondInternalError(c, err)
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
			respondInternalError(c, err)
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

	kugouAdapter, err := getKugouAdapterForUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load credentials"})
		return
	}

	var playURL string

	if platform == "kugou" {
		if kugouAdapter == nil || !kugouAdapter.IsAvailable() {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "kugou adapter not available"})
			return
		}
		playURL, err = kugouAdapter.GetPlayURL(songID, quality)
		if err != nil {
			respondInternalError(c, err)
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
