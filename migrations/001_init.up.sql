CREATE SCHEMA IF NOT EXISTS musicbox;

CREATE TABLE musicbox.playlists (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name TEXT NOT NULL DEFAULT '我的歌单',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE musicbox.playlist_songs (
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

CREATE TABLE musicbox.user_credentials (
    user_id BIGINT PRIMARY KEY,
    kugou_cookie BYTEA,
    kugou_token BYTEA,
    netease_cookie BYTEA,
    qq_cookie BYTEA,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
