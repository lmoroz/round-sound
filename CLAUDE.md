# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

- **Run in Development Mode**: `wails dev` (starts backend and frontend with hot reload)
- **Build Production Binary**: `wails build` (outputs to `build/bin/`)
- **Install Frontend Deps**: `cd frontend && npm install`
- **Frontend Lint**: `cd frontend && npm run lint`
- **Frontend Type Check**: `cd frontend && vue-tsc --noEmit`
- **Go Mod Tidy**: `go mod tidy`
- **Install Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

## Architecture & Structure

This is a **Wails v2** application combining a **Go** backend with a **Vue 3/TypeScript** frontend.

### Backend (`/` root)
- **Entry Point**: `main.go` initializes the Wails application.
- **App Logic** (`app/`):
  - `app.go`: Main application lifecycle and logic.
  - `window_windows.go`: Handles Windows-specific window management, specifically keeping the widget at desktop level (`HWND_BOTTOM`) using `SetWindowPos`.
  - `config.go`: Configuration management.
- **Media & Audio** (`media/`):
  - `audiolevels.go`: Captures system audio using **WASAPI** loopback and performs real-time FFT analysis.
  - `webnowplaying.go`: Runs a WebSocket server to receive metadata from the WebNowPlaying browser extension.

### Frontend (`frontend/`)
- **Framework**: Vue 3 + Vite + TypeScript.
- **Components**: Located in `src/components/`.
  - `CircularWidget.vue`: The main container.
  - `AudioLevelsRays.vue`: Renders the audio visualization using Canvas.
  - `MediaControls.vue`, `TrackInfo.vue`: UI for player interaction.
- **State**: Uses composables (`useMediaPlayer.ts`, `useAudioLevels.ts`) to manage state synced from the backend.

### Key Features & patterns
- **Audio Visualization**: Backend calculates FFT data and emits events to frontend (~60fps) for rendering.
- **WebNowPlaying**: Connects to browser extensions via WebSocket (default port 8974). Handles port conflicts (e.g., with Rainmeter) by allowing custom ports.
- **Window Management**: The app is designed to be a "widget" that stays behind other windows. It periodically checks Z-order.
