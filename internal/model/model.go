package model

type Song struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Artist        string   `json:"artist"`
	Album         string   `json:"album"`
	Duration      int      `json:"duration"`
	Platform      string   `json:"platform"`
	PlatformID    string   `json:"platform_id"`
	Qualities     []string `json:"qualities"`
}

type Playlist struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PlaylistSong struct {
	ID             int64  `json:"id"`
	PlaylistID     int64  `json:"playlist_id"`
	Title          string `json:"title"`
	Artist         string `json:"artist"`
	Album          string `json:"album"`
	Duration       int    `json:"duration"`
	Platform       string `json:"platform"`
	PlatformSongID string `json:"platform_song_id"`
	PlayURL        string `json:"play_url"`
	Quality        string `json:"quality"`
	SortOrder      int    `json:"sort_order"`
	CreatedAt      string `json:"created_at"`
}

type Credentials struct {
	UserID        int64  `json:"user_id"`
	KugouCookie   string `json:"kugou_cookie,omitempty"`
	KugouToken    string `json:"kugou_token,omitempty"`
	NeteaseCookie string `json:"netease_cookie,omitempty"`
	QQCookie      string `json:"qq_cookie,omitempty"`
}

type CredentialsStatus struct {
	Kugou   bool `json:"kugou"`
	Netease bool `json:"netease"`
	QQ      bool `json:"qq"`
}
