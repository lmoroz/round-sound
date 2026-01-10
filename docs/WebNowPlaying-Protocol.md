# WebNowPlaying Protocol Documentation

Документация протокола WebSocket для WebNowPlaying адаптеров.

## Версии протокола (Communication Revision)

### Handshake после подключения

Адаптер **ОБЯЗАН** отправить версию сразу после подключения:
```
ADAPTER_VERSION <version>;WNPLIB_REVISION <revision>
```

Пример:
```
ADAPTER_VERSION 1.0.0;WNPLIB_REVISION 3
```

### Доступные ревизии

- `legacy` — WebNowPlaying for Rainmeter < 0.5.0 (устаревший)
- `1` — WebNowPlaying v1.0.0-v1.0.5
- `2` — WebNowPlaying v1.1.0+
- **`3`** — **Текущий стандарт** (используем мы)

---

## Revision 3 Protocol

### Типы сообщений (MessageType)

```typescript
enum MessageType {
  PLAYER_ADDED = 0,
  PLAYER_UPDATED = 1,
  PLAYER_REMOVED = 2,
  EVENT_RESULT = 3,
  USE_DESKTOP_PLAYERS = 4,
}
```

### Формат данных Player

Pipe-separated значения (26 полей):
```
<id>|<name>|<title>|<artist>|<album>|<cover>|<state>|<position>|<duration>|<volume>|<rating>|<repeat>|<shuffle>|<ratingSystem>|<availableRepeat>|<canSetState>|<canSkipPrevious>|<canSkipNext>|<canSetPosition>|<canSetVolume>|<canSetRating>|<canSetRepeat>|<canSetShuffle>|<createdAt>|<updatedAt>|<activeAt>|
```

Примечания:
- Пустые строки заменяются на `\x01` (ASCII 1)
- Pipe символы экранируются: `|` → `\\|`
- Boolean → `1` (true) / `0` (false)

---

## Входящие сообщения (от браузера к адаптеру)

### 1. PLAYER_ADDED

Формат:
```
0 <playerId> <playerData>
```

Пример:
```
0 723877693416 |Yandex Music|Crashing Down|Army Of Lovers|...|
```

### 2. PLAYER_UPDATED (Partial Update)

Формат:
```
1 <playerId> <partialPlayerData>
```

Пример (обновление только position):
```
1 723877693416 |||||||145|||||||||||||||||1768287902026||
```

### 3. PLAYER_REMOVED

Формат:
```
2 <playerId>
```

### 4. EVENT_RESULT

Формат:
```
3 <playerId> <eventId> <statusCode>
```

**Status Codes:**
- `0` = Success
- `1` = Not Supported
- `2` = Timeout / Unable to execute

Пример:
```
3 723877693416 abc123 0
```

### 5. Бинарное сообщение (обложка)

Структура:
```
[4 bytes: playerId (UInt32 LE)][PNG data]
```

---

## Исходящие сообщения (от адаптера к браузеру)

### Формат команды (Revision 3)

```
<playerId> <eventId> <eventType> <data>
```

Где:
- `playerId` — ID плеера (число)
- `eventId` — уникальный ID команды для трекинга результата (строка)
- `eventType` — **числовой код** команды
- `data` — данные команды (опционально)

### Event Types (команды)

```typescript
enum Events {
  TRY_SET_STATE = 0,       // data: StateMode (0=STOPPED, 1=PLAYING, 2=PAUSED)
  TRY_SKIP_PREVIOUS = 1,   // data: нет
  TRY_SKIP_NEXT = 2,       // data: нет
  TRY_SET_POSITION = 3,    // data: position in seconds (int)
  TRY_SET_VOLUME = 4,      // data: volume 0-100 (int)
  TRY_SET_RATING = 5,      // data: rating (0, 1-5)
  TRY_SET_REPEAT = 6,      // data: "NONE" | "ALL" | "ONE"
  TRY_SET_SHUFFLE = 7,     // data: 0 или 1
}
```

### StateMode

```typescript
enum StateMode {
  STOPPED = 0,
  PLAYING = 1,
  PAUSED = 2,
}
```

### RepeatMode

```typescript
enum Repeat {
  NONE = 1,
  ALL = 2,
  ONE = 4,
}
```

### Примеры команд

#### Play (state = PLAYING)
```
723877693416 evt_001 0 1
```

#### Pause (state = PAUSED)  
```
723877693416 evt_002 0 2
```

#### Next track
```
723877693416 evt_003 2
```

#### Previous track
```
723877693416 evt_004 1
```

#### Toggle shuffle (ON)
```
723877693416 evt_005 7 1
```

#### Set repeat mode (ALL)
```
723877693416 evt_006 6 ALL
```

#### Set rating (Like = 5)
```
723877693416 evt_007 5 5
```

#### Seek to position (120 seconds)
```
723877693416 evt_008 3 120
```

---

## Ответы (EVENT_RESULT)

После получения команды **браузерный плагин отправит результат**:

```
3 <playerId> <eventId> <statusCode>
```

**Пример успешного выполнения:**
```
3 723877693416 evt_001 0
```

**Пример ошибки (команда не поддерживается):**
```
3 723877693416 evt_002 1
```

---

## Player Capabilities

Перед отправкой команды проверяй флаги:

- `canSetState` — можно ли управлять Play/Pause
- `canSkipNext` — можно ли переключить на следующий трек
- `canSkipPrevious` — можно ли переключить на предыдущий трек
- `canSetPosition` — можно ли перемотать
- `canSetVolume` — можно ли изменить громкость
- `canSetRating` — можно ли поставить лайк/дизлайк
- `canSetRepeat` — можно ли изменить режим повтора
- `canSetShuffle` — можно ли включить shuffle

---

## Важные замечания

1. **Player ID** — это число определяется браузером, может быть очень большим (64-bit)
2. **Event ID** — генерируется адаптером, должен быть уникальным для каждой команды
3. **Все команды асинхронные** — результат придет в EVENT_RESULT
4. **Обложки** — отправляются отдельным бинарным сообщением с Player ID в первых 4 байтах
5. **Partial Updates** — браузер отправляет только измененные поля (пустые = `\x01`)

---

## Пример потока

```
→ (адаптер подключается)
← ADAPTER_VERSION 1.0.0;WNPLIB_REVISION 3

← (браузер отправляет PLAYER_ADDED)
← 0 723877693416 |Yandex Music|Song Title|Artist|...|

← (браузер отправляет обложку бинарно)
← [4 bytes: player ID][PNG data]

→ (адаптер отправляет команду PAUSE)
→ 723877693416 cmd_pause 0 2

← (браузер отвечает результатом)
← 3 723877693416 cmd_pause 0

← (браузер отправляет PLAYER_UPDATED с новым state)
← 1 723877693416 |||||||2|||||||||||||||||1768287903000||
```

---

## Useful Links

- WebNowPlaying Browser Extension: https://github.com/keifufu/WebNowPlaying
- WebNowPlaying Adapters: https://wnp.keifufu.dev/
