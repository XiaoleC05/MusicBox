package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/XiaoleC05/MusicBox/internal/config"
	"github.com/XiaoleC05/MusicBox/internal/db"
	"github.com/XiaoleC05/MusicBox/internal/handler"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func corsOrigins() []string {
	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		return strings.Split(v, ",")
	}
	return []string{"http://localhost:5173"}
}

func main() {
	config.Load()

	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	runMigrations()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-Id", "X-Username", "X-Role"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	playlistHandler := handler.NewPlaylistHandler()
	credentialsHandler := handler.NewCredentialsHandler()

	r.GET("/api/health", handler.Health)

	api := r.Group("/api")
	api.Use(handler.AuthMiddleware())
	{
		api.GET("/search", handler.Search)
		api.GET("/play/:platform/:songId", handler.Play)

		api.GET("/playlists", playlistHandler.List)
		api.POST("/playlists", playlistHandler.Create)
		api.DELETE("/playlists/:id", playlistHandler.Delete)
		api.POST("/playlists/:id/songs", playlistHandler.AddSong)
		api.DELETE("/playlists/:id/songs/:songId", playlistHandler.RemoveSong)
		api.GET("/playlists/:id/songs", playlistHandler.ListSongs)

		api.PUT("/credentials", credentialsHandler.Update)
		api.GET("/credentials/status", credentialsHandler.Status)
	}

	srv := &http.Server{
		Addr:    ":" + config.Cfg.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func runMigrations() {
	ctx := context.Background()

	migrationSQL := `
		CREATE SCHEMA IF NOT EXISTS musicbox;

		CREATE TABLE IF NOT EXISTS musicbox.playlists (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			name TEXT NOT NULL DEFAULT '我的歌单',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS musicbox.playlist_songs (
			id BIGSERIAL PRIMARY KEY,
			playlist_id BIGINT NOT NULL REFERENCES musicbox.playlists(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			artist TEXT DEFAULT '',
			album TEXT DEFAULT '',
			duration INT DEFAULT 0,
			platform TEXT NOT NULL,
			platform_song_id TEXT NOT NULL,
			play_url TEXT DEFAULT '',
			quality TEXT DEFAULT 'standard',
			sort_order INT DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS musicbox.user_credentials (
			user_id BIGINT PRIMARY KEY,
			kugou_cookie BYTEA,
			kugou_token BYTEA,
			netease_cookie BYTEA,
			qq_cookie BYTEA,
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
	`

	_, err := db.Pool.Exec(ctx, migrationSQL)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed")
}
