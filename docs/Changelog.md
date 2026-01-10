# Changelog

## [0.1.0] 2026-01-13 17:12

- Implemented settings system with localStorage persistence
- Added SettingsPanel component with modal UI
- Created dynamic color scheme generator from single primary color
- Added FFT size configuration (1024/2048/4096/8192)
- Added frequency range settings (min/max Hz)
- Integrated color utilities (HEX/RGB/HSL conversion)
- Connected settings to backend via Wails Events
- Updated AudioLevelsRays to use dynamic colors from settings
- All settings auto-save and restore between app restarts

## [0.1.0] 2026-01-13 15:54

- Changed color scheme from teal-cyan to orange tones
- Updated primary colors: #ff8c42, #ff6b35, #ffaa66
- Updated audio visualization rays gradient to match orange theme
- Increased ray length from 30px to 90px (3x) for better visibility
- Increased ray thickness from 2px to 6px (3x) for enhanced visual presence

## [0.1.0] 2026-01-13 15:39

- Implemented WASAPI audio level capture via `IAudioMeterInformation` (media/audiolevels.go)
- Real-time audio visualization with 64 frequency bands at 60 FPS
- Audio gain boost (10x) for better visual representation
- Bass boost and organic variation for frequency distribution
- Integration with frontend via Wails Events (`audio:levels`)
- COM thread-safe initialization in capture loop
- Updated documentation (README.md, todo.md)

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
