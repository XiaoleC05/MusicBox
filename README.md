# MusicBox — Cross-Platform Music Aggregator

> One search box. All your music platforms. Seamless playback.

## Why MusicBox?

Different streaming platforms have different catalogs. That song you want might be on NetEase but not on QQ Music. You end up switching between apps, managing separate playlists, and never having everything in one place.

**MusicBox** connects to the music platforms you already have accounts on. Search once across all of them, and it automatically picks the platform that has the song. One player, one playlist, zero app switching.

## How It Works

MusicBox doesn't provide music — it uses your existing platform accounts:

1. You configure your login credentials (cookies) for each platform
2. When you search, MusicBox queries all connected platforms simultaneously
3. It picks the best available source — highest quality, fastest response
4. If the current platform can't play a track, it silently switches to another

**All playback uses your own accounts and memberships.** MusicBox is a player, not a piracy tool.

## Features

| Feature | What You Get |
|---------|-------------|
| **Multi-Platform Search** | Search by song name or artist across all connected platforms at once |
| **Auto Source Switching** | Song not available on Platform A? Instantly falls back to Platform B |
| **Full Player** | Play, pause, skip, seek, volume — everything you expect |
| **Cross-Platform Playlists** | A single playlist can contain songs from different platforms |
| **High Quality** | Prefers lossless / high-bitrate sources when available |
| **Persistent Login** | Configure once, stays logged in across sessions |

## MVP Scope

First platform: **Kugou Music**. Additional platforms will be added after the integration pattern is validated.

## Tech Stack

| Environment | Backend | Database | Frontend | Special |
|-------------|---------|----------|----------|---------|
| Online (Personal) | Go | PostgreSQL | React (Oxelia51) | Platform HTTP adapters |
| Desktop (exe) | Go | SQLite | Embedded React | Same, packaged as exe |

- **Plugin architecture**: Each music platform has its own adapter, making it easy to add new platforms
- **Online version**: Personal use only — not available to other platform users
- **Desktop version**: Full-featured exe for anyone

## Getting Started

### Desktop (exe) — Recommended for Everyone

1. Download `MusicBox.exe` from [GitHub Releases](https://github.com/XiaoleC05/MusicBox/releases)
2. Run the executable
3. Configure your music platform cookies in settings
4. Start searching and playing

### Online — Developer Only

The online version on [oxelia51.com](https://oxelia51.com) is for the developer's personal use only. It is not available to other users. Please use the exe version.

## Legal

- Uses only your own legally obtained accounts and memberships
- Does not bypass DRM, crack, or distribute copyrighted content
- You are responsible for complying with each platform's terms of service

## Status

Concept phase. Development not yet started.
