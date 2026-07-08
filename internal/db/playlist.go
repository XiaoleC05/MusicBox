package db

import (
	"context"
	"time"

	"github.com/XiaoleC05/MusicBox/internal/model"
	"github.com/jackc/pgx/v5"
)

type PlaylistRepository struct{}

func NewPlaylistRepository() *PlaylistRepository {
	return &PlaylistRepository{}
}

func (r *PlaylistRepository) ListByUser(ctx context.Context, userID int64) ([]model.Playlist, error) {
	query := `SELECT id, user_id, name, created_at, updated_at 
			  FROM musicbox.playlists 
			  WHERE user_id = $1 
			  ORDER BY created_at DESC`

	rows, err := Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []model.Playlist
	for rows.Next() {
		var p model.Playlist
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		p.CreatedAt = createdAt.Format(time.RFC3339)
		p.UpdatedAt = updatedAt.Format(time.RFC3339)
		playlists = append(playlists, p)
	}

	return playlists, rows.Err()
}

func (r *PlaylistRepository) Create(ctx context.Context, userID int64, name string) (*model.Playlist, error) {
	query := `INSERT INTO musicbox.playlists (user_id, name) 
			  VALUES ($1, $2) 
			  RETURNING id, user_id, name, created_at, updated_at`

	var p model.Playlist
	var createdAt, updatedAt time.Time
	err := Pool.QueryRow(ctx, query, userID, name).Scan(
		&p.ID, &p.UserID, &p.Name, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.CreatedAt = createdAt.Format(time.RFC3339)
	p.UpdatedAt = updatedAt.Format(time.RFC3339)
	return &p, nil
}

func (r *PlaylistRepository) Delete(ctx context.Context, playlistID, userID int64) error {
	query := `DELETE FROM musicbox.playlists WHERE id = $1 AND user_id = $2`
	_, err := Pool.Exec(ctx, query, playlistID, userID)
	return err
}

// VerifyPlaylistOwner checks if the playlist belongs to the user.
func (r *PlaylistRepository) VerifyPlaylistOwner(ctx context.Context, playlistID, userID int64) (bool, error) {
	var exists bool
	err := Pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM musicbox.playlists WHERE id = $1 AND user_id = $2)`,
		playlistID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *PlaylistRepository) AddSong(ctx context.Context, playlistID, userID int64, song model.PlaylistSong) error {
	// Verify ownership first
	owned, err := r.VerifyPlaylistOwner(ctx, playlistID, userID)
	if err != nil {
		return err
	}
	if !owned {
		return pgx.ErrNoRows
	}

	query := `INSERT INTO musicbox.playlist_songs 
			  (playlist_id, title, artist, album, duration, platform, platform_song_id, play_url, quality, sort_order) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = Pool.Exec(ctx, query,
		playlistID, song.Title, song.Artist, song.Album, song.Duration,
		song.Platform, song.PlatformSongID, song.PlayURL, song.Quality, song.SortOrder,
	)
	return err
}

func (r *PlaylistRepository) RemoveSong(ctx context.Context, songID, userID int64) error {
	query := `DELETE FROM musicbox.playlist_songs 
			  WHERE id = $1 
			  AND playlist_id IN (SELECT id FROM musicbox.playlists WHERE user_id = $2)`
	_, err := Pool.Exec(ctx, query, songID, userID)
	return err
}

func (r *PlaylistRepository) ListSongs(ctx context.Context, playlistID, userID int64) ([]model.PlaylistSong, error) {
	query := `SELECT ps.id, ps.playlist_id, ps.title, ps.artist, ps.album, ps.duration, ps.platform, 
			  ps.platform_song_id, ps.play_url, ps.quality, ps.sort_order, ps.created_at 
			  FROM musicbox.playlist_songs ps
			  JOIN musicbox.playlists p ON p.id = ps.playlist_id
			  WHERE ps.playlist_id = $1 AND p.user_id = $2
			  ORDER BY ps.sort_order, ps.created_at`

	rows, err := Pool.Query(ctx, query, playlistID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []model.PlaylistSong
	for rows.Next() {
		var s model.PlaylistSong
		var createdAt time.Time
		if err := rows.Scan(&s.ID, &s.PlaylistID, &s.Title, &s.Artist, &s.Album,
			&s.Duration, &s.Platform, &s.PlatformSongID, &s.PlayURL, &s.Quality,
			&s.SortOrder, &createdAt); err != nil {
			return nil, err
		}
		s.CreatedAt = createdAt.Format(time.RFC3339)
		songs = append(songs, s)
	}

	return songs, rows.Err()
}
