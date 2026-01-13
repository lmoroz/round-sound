# Changelog

## [0.3.1] 2026-01-14 01:50

- **Interactive progress bar thumb**: Added draggable thumb button that appears on hover
- Circular thumb follows progress arc with glow effects and inner highlight
- Drag-to-seek functionality for track position control
- Click-on-track to seek to specific position
- ViewBox padding to prevent thumb clipping at edges
- Added `--wails-draggable: no-drag` to prevent window dragging during seek
- Touch support for touchscreen devices
- Fixed nil pointer dereference when last player is removed
- Rounded seek position to integer before sending to backend

## [0.3.0] 2026-01-13 19:15

- **Window resize**: Increased application window from 400x400px to 600x600px
- **Audio visualization expansion**: Canvas size increased to 580px to fully display rays without clipping
- **System tray integration**: Using getlantern/systray with Register() for non-blocking tray icon
- **Taskbar hiding**: Window hidden from taskbar using WS_EX_TOOLWINDOW style
- **Tray menu**: Added "Выход" menu item for application exit
- **Autorun functionality**: Added Windows registry-based autorun manager (app/autorun.go)
- **Settings UI enhancement**: Added "System" section with autorun toggle in SettingsPanel
- **Application icon**: Custom icon embedded via go:embed (tray.ico)
- **Go backend extensions**: Added TrayManager with systray.Register() integration
- **UI adjustments**: Slightly increased cover size to 290px and reduced background opacity to 0.55
- Updated documentation (README.md, todo.md) with new features

## [0.2.0] 2026-01-13 17:XX

- **Implemented FFT-based audio analysis** replacing peak-distribution visualization
- Switched from `IAudioMeterInformation` to WASAPI loopback capture (`IAudioCaptureClient`) for raw audio samples
- Integrated `github.com/mjibson/go-dsp/fft` for real-time FFT processing
- Added Hann window function to reduce spectral leakage
- Implemented logarithmic frequency binning (20Hz - 20kHz) into 64 bands
- Dynamic FFT configuration: adjustable FFT size (1024/2048/4096/8192) and frequency range
- Frontend audio config event handling (`audio:config`) for real-time settings updates
- New module `media/fft.go` for FFT processing logic
- Complete rewrite of `media/audiolevels.go` for loopback capture and FFT integration
- Updated documentation (README.md, todo.md) to reflect FFT implementation

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
