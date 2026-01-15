# Round Sound Widget — AI Assistant Guide

**Версия проекта**: 0.3.4 (на 2026-01-15)
**Цель проекта**: Круговой музыкальный виджет для Windows с визуализацией аудио-уровней в реальном времени

---

## 🎯 ЧТО ДЕЛАТЬ ПЕРЕД НАЧАЛОМ РАБОТЫ

1. **Прочитай этот файл целиком** — здесь вся ключевая информация о проекте
2. **Изучи `docs/Changelog.md`** — узнай, что было реализовано недавно
3. **Проверь `docs/todo.md`** — планы разработки и известные проблемы
4. **Посмотри `README.md`** — общее описание и features
5. **НЕ РЕДАКТИРУЙ** файлы из `.gitignore` (особенно `WebNowPlaying/` и `AudioLevel/` — это reference-проекты)

---

## 📂 СТРУКТУРА ПРОЕКТА

### Backend (Go + Wails v2)

```text
/
├── main.go                  # Entry point: конфигурация Wails, окно 600x600px, WebView2
├── wails.json               # Конфигурация Wails (версия: 0.3.3, команды сборки)
├── go.mod / go.sum          # Go dependencies
├── app/                     # Core application logic
│   ├── app.go               # App struct, Startup/Shutdown/DomReady, SaveWindowPosition
│   ├── config.go            # Config persistence (%APPDATA%/round-sound/config.json)
│   ├── window.go            # Cross-platform window manager
│   ├── window_windows.go    # Windows-specific: HWND_BOTTOM (desktop level)
│   ├── tray.go              # System tray integration (getlantern/systray)
│   └── autorun.go           # Windows registry autorun manager
├── media/                   # Media player & audio capture
│   ├── types.go             # PlayerData, Rating, RepeatMode, etc.
│   ├── webnowplaying.go     # WebSocket server (default port 8974), WNP protocol rev.3
│   ├── audiolevels.go       # WASAPI loopback capture (IAudioCaptureClient)
│   └── fft.go               # FFT processing (Hann window, 64 bands, 20Hz-20kHz)
└── build/                   # Build artifacts (wails build output)
```

**Ключевые Go модули**:

- **Wails v2**: desktop framework (embed, runtime, events)
- **gorilla/websocket**: WebNowPlaying communication
- **go-wca**: WASAPI для захвата аудио
- **mjibson/go-dsp/fft**: FFT анализ
- **getlantern/systray**: system tray

### Frontend (Vue 3 + TypeScript + Vite 7)

```text
frontend/
├── package.json             # Dependencies: Vue 3.5.26, Vite 7.3.1, lucide-vue-next
├── vite.config.ts           # Vite config: alias @, checker, eslint (dev only)
├── tsconfig.json            # TypeScript config
├── eslint.config.js         # ESLint config (flat config)
├── src/
│   ├── main.ts              # Entry point
│   ├── App.vue              # Root component (содержит CircularWidget)
│   ├── components/
│   │   ├── CircularWidget.vue       # Main widget: position, drag, settings button
│   │   ├── AlbumCover.vue           # Album art с fallback
│   │   ├── TrackInfo.vue            # Title + Artist
│   │   ├── ProgressRing.vue         # SVG progress ring с draggable thumb
│   │   ├── AudioLevelsRays.vue      # Canvas FFT visualization (64 bars)
│   │   ├── MediaControls.vue        # Play/Pause/Next/Prev/Shuffle/Repeat/Like
│   │   ├── SettingsPanel.vue        # Modal settings (colors, FFT, autorun, WNP port)
│   │   └── ContextMenu.vue          # Right-click context menu (Выход)
│   ├── composables/
│   │   ├── useApp.ts                # App lifecycle (CloseApp IPC binding)
│   │   ├── useMediaPlayer.ts        # WNP state, EventsOn('player:update'), commands
│   │   ├── useAudioLevels.ts        # EventsOn('audio:levels'), canvas rendering
│   │   └── useSettings.ts           # Settings state, localStorage, EventsEmit
│   ├── types/
│   │   ├── index.ts                 # PlayerData, Settings, RepeatMode, Rating
│   │   └── wails.d.ts               # Centralized Wails runtime types
│   └── utils/
│       └── colors.ts                # HEX/RGB/HSL conversion, color scheme generation
└── wailsjs/                         # Auto-generated Go bindings (не редактировать вручную!)
    ├── go/app/App.js                # Go method bindings
    └── runtime/                     # Wails runtime API
```

**Ключевые npm packages**:

- **Vue 3.5.26** (Composition API, `<script setup>`)
- **Vite 7.3.1** (dev server, build)
- **lucide-vue-next 0.562.0** (иконки)
- **wire-ts / wire-vue** (DI framework, используется для composables)
- **ESLint + vue-tsc** (linting & type checking)

### Reference Projects (НЕ РЕДАКТИРОВАТЬ!)

- **`WebNowPlaying/`** — исходники браузерного плагина WebNowPlaying (git submodule/reference)
- **`AudioLevel/`** — reference/example для WASAPI audio capture

---

## 🔑 КЛЮЧЕВЫЕ КОНЦЕПЦИИ

### 1. АРХИТЕКТУРА

```text
┌─────────────────────────────────────────────────────────────────┐
│                          Browser (WebNowPlaying Plugin)          │
│                     (YouTube Music / Spotify Web)                │
└────────────────────────────┬────────────────────────────────────┘
                             │ WebSocket (port 8974 or custom)
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Backend (Go + Wails)                          │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ media/webnowplaying.go                                    │  │
│  │  - WebSocket server                                       │  │
│  │  - Protocol rev.3 (partial updates, commands)             │  │
│  │  - EventsEmit('player:update', PlayerData)                │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ media/audiolevels.go + fft.go                             │  │
│  │  - WASAPI loopback capture                                │  │
│  │  - FFT analysis (64 bands, Hann window)                   │  │
│  │  - EventsEmit('audio:levels', []float64)                  │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ app/app.go, config.go, window.go, tray.go, autorun.go    │  │
│  │  - Window position persistence                            │  │
│  │  - Desktop-level window (HWND_BOTTOM)                     │  │
│  │  - System tray, autorun                                   │  │
│  └───────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             │ Wails Events + Bindings
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Frontend (Vue 3 + TypeScript)                │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ composables/useMediaPlayer.ts                             │  │
│  │  - EventsOn('player:update') → reactive state             │  │
│  │  - Commands: Play, Pause, Next, Prev, SetPosition, etc.  │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ composables/useAudioLevels.ts                             │  │
│  │  - EventsOn('audio:levels') → canvas rendering            │  │
│  │  - DPR change detection & auto-recovery                   │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ components/CircularWidget.vue                             │  │
│  │  - Draggable position, volume control (wheel)             │  │
│  │  - Context menu (right-click)                             │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### 2. WAILS EVENTS (Backend → Frontend)

| Event Name      | Payload Type     | Purpose                                                           |
|-----------------|------------------|-------------------------------------------------------------------|
| `player:update` | `PlayerData`     | Обновление данных плеера (title, artist, position, etc.)         |
| `audio:levels`  | `[]float64` (64) | FFT данные для визуализации (60 FPS)                              |
| `audio:config`  | `AudioConfig`    | Обновление FFT конфигурации из settings                           |

### 3. WAILS BINDINGS (Frontend → Backend, Go methods)

**App methods** (`app/app.go`):

- `SaveWindowPosition(x, y int)`
- `LoadWindowPosition() (int, int)`
- `GetVersion() string`
- `CloseApp()` — graceful shutdown
- `SetAutorun(enabled bool)`
- `IsAutorunEnabled() bool`
- `ChangeWNPPort(port int)`
- `GetWNPPort() int`
- `IsWNPConnected() bool`

**Media methods** (через `useMediaPlayer.ts`):

- `Play()`, `Pause()`, `PlayPause()`
- `Next()`, `Previous()`
- `SetPosition(seconds int)`
- `SetShuffle(enabled bool)`
- `SetRepeat(mode string)` — "NONE" / "ONE" / "ALL"
- `ToggleLike()`, `ToggleDislike()`
- `SetVolume(percent int)` — 0-100

### 4. SETTINGS PERSISTENCE

- **Backend**: `%APPDATA%/round-sound/config.json` (Go `app/config.go`)
  - Window position (X, Y)
  - WNP port

- **Frontend**: `localStorage` (`useSettings.ts`)
  - Primary color
  - FFT size (1024/2048/4096/8192)
  - Frequency range (minHz, maxHz)
  - Autorun toggle (синхронизируется с backend)

### 5. ВАЖНЫЕ ОСОБЕННОСТИ КОДА

#### Frontend Code Style (user rules)

- **Отступы**: 2 пробела
- **Условия**:
  - Одиночный statement → без `{}`, одна строка: `if (x) doSomething();`
  - Множественные statements → обязательно `{}`
  - `else` всегда на новой строке
- **Импорты**:
  - ❌ Относительные пути `../components/Foo.vue`
  - ✅ Alias `@/components/Foo.vue`
  - ❌ Ручной импорт Vue-компонентов (есть `unplugin-vue-components`)
  - ✅ Автоимпорт компонентов по имени
- **Vue файлы**: `<script setup>` **ВСЕГДА ВЫШЕ** `<template>`

#### Go Code

- COM thread safety: `runtime.LockOSThread()` в WASAPI loops
- Mutex locks: все методы `media/webnowplaying.go` защищены от race conditions
- Error handling: подробный logging в консоль

---

## 🛠️ РАБОЧИЙ ПРОЦЕСС (WORKFLOW)

### Development

```bash
# Первичная установка зависимостей
cd frontend && npm install && cd ..
go mod tidy

# Запуск dev mode (hot reload frontend + backend)
wails dev
```

### Build Production

```bash
wails build
# Output: build/bin/round-sound.exe (Windows)
```

### Linting & Type Check

```bash
cd frontend
npm run lint         # ESLint + vue-tsc
npm run lint:fix     # Auto-fix
```

---

## 📝 CHANGELOG & TODO

- **Изменения**: см. `docs/Changelog.md` (версии с 0.1.0 до 0.3.4)
- **Планы**: см. `docs/todo.md` (UI улучшения, multi-player support, installer)

---

## ⚠️ ЧАСТЫЕ ОШИБКИ И ЛОВУШКИ

### 1. WebNowPlaying Port Conflict

- **Проблема**: Порт 8974 занят (Rainmeter + WebNowPlaying.dll)
- **Решение**: Custom Adapter (настройка в SettingsPanel, сохранение в config.json)

### 2. Audio Levels Не Отображаются

- **Причины**:
  - WASAPI не инициализирован (COM error)
  - Canvas DPR mismatch (после смены масштаба экрана)
- **Решение**:
  - Проверь логи в консоли Go
  - `useAudioLevels.ts` автоматически переинициализирует canvas при изменении DPR

### 3. Window Не На Desktop Level

- **Причина**: Z-order сброшен (Windows API)
- **Решение**: `window_windows.go` автоматически вызывает `SetWindowPos(HWND_BOTTOM)` каждые 500ms

### 4. Wails Bindings Не Обновляются

- **Решение**:
  - `wails dev` авто-генерирует `frontend/wailsjs/`
  - НЕ редактируй `wailsjs/` вручную!
  - Restart `wails dev` если bindings сломались

### 5. ESLint/TypeScript Ошибки

- **Решение**:
  - `npm run lint:fix` — авто-фикс
  - Проверь `.editorconfig`, `eslint.config.js`, `tsconfig.json`

---

## 📚 РЕФЕРЕНС-ПРОЕКТЫ (Read-Only)

- **`./WebNowPlaying`** — исходники браузерного плагина (protocol reference)
- **`./AudioLevel`** — пример WASAPI loopback capture (Go reference)

❌ **НЕ РЕДАКТИРУЙ ЭТИ ПАПКИ** — они в `.gitignore` и используются только для справки!

---

## 🎨 UI/UX ОСОБЕННОСТИ

### Визуализация FFT Rays

- 64 полосы частот (20Hz - 20kHz, логарифмическая шкала)
- Окраска: градиент orange → золото (если ANY аудио > 0.02), серый (тишина)
- Canvas: 580x580px (window 600x600px - padding)
- DPR-aware: автоматический resize при изменении масштаба экрана

### Progress Ring с Draggable Thumb

- SVG arc + circular thumb button
- Появляется при hover
- Drag-to-seek + click-to-seek
- `--wails-draggable: no-drag` — предотвращает перетаскивание окна во время seek

### Volume Control

- Mouse wheel over widget → EventsEmit('media:volume')
- Floating overlay с процентами и иконкой
- Volume indicator внизу виджета при hover

### Context Menu

- Right-click → ContextMenu.vue
- "Выход" → CloseApp() → graceful shutdown (tray cleanup, save config)

---

## 🚀 ГОТОВ К РАБОТЕ

**Перед началом задачи**:

1. Прочитай user request
2. Определи, какие файлы нужно изменить (backend / frontend / оба)
3. Проверь историю изменений в Changelog.md
4. Следуй code style rules (user_rules)
5. Тестируй изменения в `wails dev`
6. Обнови `docs/Changelog.md` если задача завершена

**Удачи!** 🎵
