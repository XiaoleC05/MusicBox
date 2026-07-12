package handler

import (
	"net/http"

	"github.com/XiaoleC05/MusicBox/internal/adapter"
	"github.com/XiaoleC05/MusicBox/internal/config"
	"github.com/XiaoleC05/MusicBox/internal/db"
	"github.com/gin-gonic/gin"
)

type CredentialsHandler struct {
	repo *db.CredentialsRepository
}

func NewCredentialsHandler() *CredentialsHandler {
	return &CredentialsHandler{
		repo: db.NewCredentialsRepository(),
	}
}

func (h *CredentialsHandler) Update(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	var req struct {
		KugouCookie   string `json:"kugou_cookie"`
		KugouToken    string `json:"kugou_token"`
		NeteaseCookie string `json:"netease_cookie"`
		QQCookie      string `json:"qq_cookie"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key := config.Cfg.EncryptionKey

	var kugouCookie, kugouToken, neteaseCookie, qqCookie string
	var err error

	if req.KugouCookie != "" {
		kugouCookie, err = adapter.Encrypt(req.KugouCookie, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption failed"})
			return
		}
	}

	if req.KugouToken != "" {
		kugouToken, err = adapter.Encrypt(req.KugouToken, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption failed"})
			return
		}
	}

	if req.NeteaseCookie != "" {
		neteaseCookie, err = adapter.Encrypt(req.NeteaseCookie, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption failed"})
			return
		}
	}

	if req.QQCookie != "" {
		qqCookie, err = adapter.Encrypt(req.QQCookie, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption failed"})
			return
		}
	}

	if err := h.repo.Upsert(c.Request.Context(), userID, kugouCookie, kugouToken, neteaseCookie, qqCookie); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "credentials updated"})
}

func (h *CredentialsHandler) Status(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	status, err := h.repo.GetStatus(c.Request.Context(), userID)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, status)
}
