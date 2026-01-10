## [0.1.0] 2026-01-13 15:12

- Implemented WebNowPlaying Revision 3 protocol for media controls
- Fixed command format: using numeric event types and proper message structure
- Added event result handling with detailed logging
- Created comprehensive protocol documentation (`docs/WebNowPlaying-Protocol.md`)
- Fixed race conditions in media control methods with proper mutex locks
- All media control buttons (play/pause, next, previous, shuffle, repeat, rating) now functional
- Added debug logging throughout frontend and backend for easier troubleshooting

## [0.1.0] 2026-01-13 14:55

- Initial project setup with Wails v2 + Vue 3 + TypeScript
- Implemented desktop-level frameless window (HWND_BOTTOM)
- Added WebNowPlaying bidirectional WebSocket integration (port 8974)
- Implemented core UI: CircularWidget, AlbumCover, TrackInfo, MediaControls (with Heart icon), ProgressRing
- Added AudioLevelsRays canvas visualization (mock data)
- Implemented window position saving and restoration
- Fixed port binding issues and added debug logging for WNP
