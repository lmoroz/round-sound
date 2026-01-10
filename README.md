# Round Sound Widget

A mini-desktop application for Windows based on Wails + Vue.js + Vite + TypeScript.

## Features
- **Widget**: Round, transparent, draggable window.
- **Integration**: Works with **WebNowPlaying** browser extension (Port 8974).
- **Controls**: Play, Pause, Next, Prev, Shuffle, Repeat, Like/Dislike.
- **Tech Stack**: Wails (Go), Vue 3, TailwindCSS, Lucide Icons.

## Setup

1. **Prerequisites**:
   - Go 1.21+
   - Node.js & pnpm
   - Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

2. **Install Dependencies**:
   ```bash
   cd frontend
   pnpm install
   cd ..
   ```

3. **Run Development**:
   ```bash
   wails dev
   ```

4. **Build for Production**:
   ```bash
   wails build
   ```

## Architecture
- **Backend (`main.go`, `pkg/wnp`)**: Handles WebSocket communication with the browser extension.
- **Frontend (`frontend/src`)**: Vue 3 application with a circular design.

## Usage
Ensure your browser has the **WebNowPlaying** extension installed and active. The widget will automatically display the playing media.
