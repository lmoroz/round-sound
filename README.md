# Round Sound Widget

Круглый музыкальный виджет для Windows с визуализацией audio levels.

## Особенности

- 🎵 Отображение информации о треке (название, артист, обложка)
- 🎮 Кнопки управления (play/pause, next, previous, shuffle, repeat, like/dislike)
- 📊 Визуализация audio levels "лучами" вокруг виджета
- 🖼️ Круговой прогресс-бар
- 🪟 Прозрачное frameless окно на уровне desktop
- 💾 Запоминание позиции окна между перезапусками
- ⚙️ Настройки с динамической цветовой схемой и параметрами FFT

## Технологии

### Backend
- **Go 1.21+**
- **Wails v2** — десктопный фреймворк
- **gorilla/websocket** — WebSocket для WebNowPlaying
- **go-wca** — WASAPI для захвата audio levels

### Frontend
- **Vue.js 3** + Composition API
- **TypeScript**
- **Vite 7**
- **Lucide Icons**
- **Canvas** для audio levels

## Структура проекта

```
round-sound/
├── main.go                 # Entry point Wails
├── wails.json              # Конфигурация Wails
├── go.mod                  # Go модуль
├── app/
│   ├── app.go              # Основная логика приложения
│   ├── config.go           # Управление конфигурацией
│   ├── window.go           # Window manager (общий)
│   └── window_windows.go   # Windows-специфичный код (HWND_BOTTOM)
├── media/
│   ├── types.go            # Типы данных Player
│   ├── webnowplaying.go    # WebSocket server для WebNowPlaying
│   └── audiolevels.go      # WASAPI audio capture
├── frontend/
│   ├── src/
│   │   ├── main.ts
│   │   ├── App.vue
│   │   ├── components/
│   │   │   ├── CircularWidget.vue
│   │   │   ├── AlbumCover.vue
│   │   │   ├── TrackInfo.vue
│   │   │   ├── ProgressRing.vue
│   │   │   ├── AudioLevelsRays.vue
│   │   │   └── MediaControls.vue
│   │   ├── composables/
│   │   │   ├── useMediaPlayer.ts
│   │   │   └── useAudioLevels.ts
│   │   └── types/
│   │       └── index.ts
│   ├── package.json
│   └── vite.config.ts
└── docs/
    ├── todo.md
    └── WebNowPlaying-Protocol.md
```

## Установка и запуск

### Требования
- Go 1.21+
- Node.js 20+
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Разработка
```bash
# Установить зависимости
cd frontend && npm install && cd ..
go mod tidy

# Запустить в dev режиме
wails dev
```

### Сборка
```bash
wails build
```

## Интеграция с WebNowPlaying

Виджет работает с браузерным плагином [WebNowPlaying](https://wnp.keifufu.dev/):
1. Установите плагин в Chrome/Firefox
2. Запустите Round Sound
3. Откройте YouTube Music, Spotify Web или другой поддерживаемый сервис

### Поддерживаемые источники
- YouTube Music
- Spotify Web
- SoundCloud
- Deezer
- Tidal
- Apple Music
- И другие...

## Особенности реализации

### Desktop-level окно
Виджет отображается на уровне рабочего стола (ниже всех окон) с помощью:
- Windows API `SetWindowPos` с `HWND_BOTTOM`
- Периодическая проверка Z-order каждые 500ms

### Partial Updates

WebNowPlaying отправляет только измененные поля. Backend хранит полное состояние в памяти и выполняет merge.

### Audio Levels (WASAPI)

Визуализация звука работает через Windows Core Audio API:

- Захват peak levels через `IAudioMeterInformation`
- Усиление 10x для лучшей видимости
- Распределение по 64 частотным полосам с bass boost
- Отправка данных во frontend ~60 FPS через Wails Events

### Обработка занятого порта

Если порт 8974 занят (например, Rainmeter), выводится уведомление.

## Лицензия

MIT
