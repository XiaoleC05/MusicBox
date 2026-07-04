# MusicBox

Cross-platform music aggregator. Connect multiple streaming accounts, search across all platforms, and play with automatic source switching.

## Features

- Search across multiple music platforms simultaneously
- Auto-switch to an alternative platform when a song is unavailable on the current one
- Full player controls: play, pause, skip, seek, volume
- Cross-platform playlists containing songs from different sources
- High-quality audio prioritization (lossless / high bitrate)
- Persistent login with cookie-based session recovery

## Architecture

```text
Browser (or embedded WebView)
  ↓
React Frontend (search + player UI)
  ↓
Go Backend (platform adapters, audio streaming)
  ├── Kugou Adapter
  ├── NetEase Adapter (planned)
  └── QQ Music Adapter (planned)

PostgreSQL / SQLite (user config, playlists)
```

The Go backend uses a plugin architecture for platform adapters. Each adapter is an independent module handling search, quality selection, and playback URL resolution. The desktop version stores user configuration and playlists in SQLite.

## Requirements

- Desktop: standalone executable, no runtime dependencies
- Login credentials (cookies) for each music platform

## Installation

### Desktop

Download `MusicBox.exe` from [GitHub Releases](https://github.com/XiaoleC05/MusicBox/releases).

### Online

The online version is for the developer's personal use only.

## Usage

### Desktop

1. Double-click `MusicBox.exe` to start
2. Enter platform cookies in settings
3. Search for songs or artists and play

### Online

The online version is for the developer's personal use only. Use the desktop version instead.

## Roadmap

- [ ] Kugou Music integration (MVP)
- [ ] NetEase Cloud Music integration
- [ ] QQ Music integration

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/xxx`)
3. Commit your changes (`git commit -m 'Add xxx'`)
4. Push the branch (`git push origin feature/xxx`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.
