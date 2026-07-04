# MusicBox

跨平台音乐聚合播放器。关联多平台账号，统一搜索，自动切换最优音源。

## Features

- 同时检索多个音乐平台，聚合搜索结果
- 当前平台无版权时自动切换备用平台
- 完整播放器控制：播放、暂停、切歌、进度、音量
- 一份歌单可包含来自不同平台的歌曲
- 高音质优先，无损/高码率自动选择
- 一次配置平台登录凭证，后续自动恢复

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

Go 后端以插件架构组织各音乐平台适配器。每个平台为独立模块，处理搜索、音质筛选和播放地址解析。桌面版使用 SQLite 存储用户配置和歌单。

## Requirements

- 桌面版：独立可执行文件，无需运行时依赖
- 各音乐平台的登录凭证（Cookie）

## Installation

### 桌面版

从 [GitHub Releases](https://github.com/XiaoleC05/MusicBox/releases) 下载 `MusicBox.exe`。

### 在线版

在线版仅供开发者个人使用。

## Usage

### 桌面

1. 双击 `MusicBox.exe` 启动
2. 在设置中为各平台填写 Cookie
3. 搜索歌曲或歌手名称，开始播放

### 在线

在线版仅供开发者个人使用，请使用桌面版。

## Roadmap

- [ ] 酷狗音乐适配（MVP）
- [ ] 网易云音乐适配
- [ ] QQ 音乐适配

## Contributing

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/xxx`)
3. 提交变更 (`git commit -m 'Add xxx'`)
4. 推送分支 (`git push origin feature/xxx`)
5. 提交 Pull Request

## License

This project is licensed under the MIT License.
