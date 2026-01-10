package media

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MessageType represents WebNowPlaying message types
type MessageType int

const (
	MessagePlayerAdded   MessageType = 0
	MessagePlayerUpdated MessageType = 1
	MessagePlayerRemoved MessageType = 2
	MessageEventResult   MessageType = 3
)

// PlayerUpdateCallback is called when player state changes
type PlayerUpdateCallback func(player *Player)

// WebNowPlayingServer manages WebSocket connections with WebNowPlaying
type WebNowPlayingServer struct {
	port           int
	server         *http.Server
	upgrader       websocket.Upgrader
	conn           *websocket.Conn
	connMu         sync.Mutex
	players        map[int]*Player
	playersMu      sync.RWMutex
	activePlayerID int
	onUpdate       PlayerUpdateCallback
	stopCh         chan struct{}
	coverDir       string
}

// NewWebNowPlayingServer creates and starts a new WebNowPlaying server
func NewWebNowPlayingServer(port int, onUpdate PlayerUpdateCallback) (*WebNowPlayingServer, error) {
	// Create cover cache directory
	coverDir := filepath.Join(os.TempDir(), "round-sound", "covers")
	os.MkdirAll(coverDir, 0755)

	s := &WebNowPlayingServer{
		port:     port,
		players:  make(map[int]*Player),
		onUpdate: onUpdate,
		stopCh:   make(chan struct{}),
		coverDir: coverDir,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for WebNowPlaying
			},
		},
	}

	// Start HTTP server for WebSocket
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleConnection)

	s.server = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		log.Printf("WebNowPlaying server starting on port %d", port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("WebNowPlaying server error: %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	return s, nil
}

// Stop stops the WebNowPlaying server
func (s *WebNowPlayingServer) Stop() {
	close(s.stopCh)
	if s.conn != nil {
		s.conn.Close()
	}
	if s.server != nil {
		s.server.Close()
	}
	log.Println("WebNowPlaying server stopped")
}

// handleConnection handles incoming WebSocket connections
func (s *WebNowPlayingServer) handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	s.connMu.Lock()
	if s.conn != nil {
		s.conn.Close()
	}
	s.conn = conn
	s.connMu.Unlock()

	log.Println("WebNowPlaying client connected")

	// Send version info
	conn.WriteMessage(websocket.TextMessage, []byte("ADAPTER_VERSION 1.0.0;WNPLIB_REVISION 3"))

	// Handle messages
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		switch messageType {
		case websocket.TextMessage:
			s.handleTextMessage(string(data))
		case websocket.BinaryMessage:
			s.handleBinaryMessage(data)
		}
	}

	s.connMu.Lock()
	if s.conn == conn {
		s.conn = nil
	}
	s.connMu.Unlock()

	log.Println("WebNowPlaying client disconnected")
}

// handleTextMessage processes text messages from WebNowPlaying
func (s *WebNowPlayingServer) handleTextMessage(msg string) {
	// Log everything for debug
	if len(msg) > 100 {
		log.Printf("Incoming WNP message (len=%d): %s...", len(msg), msg[:100])
	} else {
		log.Printf("Incoming WNP message: %s", msg)
	}

	parts := strings.SplitN(msg, " ", 3)
	if len(parts) < 2 {
		return
	}

	msgType, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Printf("Invalid message type: %s", parts[0])
		return
	}

	playerID, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("Invalid player ID: %s", parts[1])
		return
	}

	switch MessageType(msgType) {
	case MessagePlayerAdded:
		if len(parts) >= 3 {
			s.handlePlayerAdded(playerID, parts[2])
		}
	case MessagePlayerUpdated:
		if len(parts) >= 3 {
			s.handlePlayerUpdated(playerID, parts[2])
		}
	case MessagePlayerRemoved:
		s.handlePlayerRemoved(playerID)
	}
}

// handleBinaryMessage processes binary messages (cover art)
func (s *WebNowPlayingServer) handleBinaryMessage(data []byte) {
	if len(data) < 4 {
		return
	}

	// First 4 bytes are player ID (little endian)
	playerID := int(binary.LittleEndian.Uint32(data[:4]))
	coverData := data[4:]

	// Save cover to file
	coverPath := filepath.Join(s.coverDir, fmt.Sprintf("%d.png", playerID))
	if err := os.WriteFile(coverPath, coverData, 0644); err != nil {
		log.Printf("Failed to save cover: %v", err)
		return
	}

	// Update player with cover path
	s.playersMu.Lock()
	if player, ok := s.players[playerID]; ok {
		player.Cover = coverPath
		player.CoverData = coverData
		s.playersMu.Unlock()
		s.notifyUpdate(player)
	} else {
		s.playersMu.Unlock()
	}

	log.Printf("Received cover for player %d (%d bytes)", playerID, len(coverData))
}

// parsePlayerData parses pipe-separated player data
func parsePlayerData(data string) map[string]string {
	fields := []string{
		"id", "name", "title", "artist", "album", "cover",
		"state", "position", "duration", "volume", "rating",
		"repeat", "shuffle", "ratingSystem", "availableRepeat",
		"canSetState", "canSkipPrevious", "canSkipNext",
		"canSetPosition", "canSetVolume", "canSetRating",
		"canSetRepeat", "canSetShuffle", "createdAt", "updatedAt", "activeAt",
	}

	result := make(map[string]string)
	parts := strings.Split(data, "|")

	for i, part := range parts {
		if i >= len(fields) {
			break
		}
		// Unescape pipe characters
		part = strings.ReplaceAll(part, "\\|", "|")
		// Skip empty marker (ASCII 1)
		if part == "\x01" {
			part = ""
		}
		result[fields[i]] = part
	}

	return result
}

// applyPlayerData applies parsed data to player
func applyPlayerData(player *Player, data map[string]string) {
	if v, ok := data["name"]; ok && v != "" {
		player.Name = v
	}
	if v, ok := data["title"]; ok && v != "" {
		player.Title = v
	}
	if v, ok := data["artist"]; ok && v != "" {
		player.Artist = v
	}
	if v, ok := data["album"]; ok && v != "" {
		player.Album = v
	}
	if v, ok := data["cover"]; ok && v != "" {
		player.Cover = v
	}
	if v, ok := data["state"]; ok && v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			player.State = StateMode(val)
		}
	}
	if v, ok := data["position"]; ok && v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			player.Position = val
		}
	}
	if v, ok := data["duration"]; ok && v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			player.Duration = val
		}
	}
	if v, ok := data["volume"]; ok && v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			player.Volume = val
		}
	}
	if v, ok := data["rating"]; ok && v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			player.Rating = val
		}
	}
	if v, ok := data["repeat"]; ok && v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			player.Repeat = RepeatMode(val)
		}
	}
	if v, ok := data["shuffle"]; ok && v != "" {
		player.Shuffle = v == "1"
	}
	if v, ok := data["ratingSystem"]; ok && v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			player.RatingSystem = RatingSystem(val)
		}
	}
	if v, ok := data["canSetState"]; ok && v != "" {
		player.CanSetState = v == "1"
	}
	if v, ok := data["canSkipPrevious"]; ok && v != "" {
		player.CanSkipPrevious = v == "1"
	}
	if v, ok := data["canSkipNext"]; ok && v != "" {
		player.CanSkipNext = v == "1"
	}
	if v, ok := data["canSetPosition"]; ok && v != "" {
		player.CanSetPosition = v == "1"
	}
	if v, ok := data["canSetVolume"]; ok && v != "" {
		player.CanSetVolume = v == "1"
	}
	if v, ok := data["canSetRating"]; ok && v != "" {
		player.CanSetRating = v == "1"
	}
	if v, ok := data["canSetRepeat"]; ok && v != "" {
		player.CanSetRepeat = v == "1"
	}
	if v, ok := data["canSetShuffle"]; ok && v != "" {
		player.CanSetShuffle = v == "1"
	}
	if v, ok := data["activeAt"]; ok && v != "" {
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			player.ActiveAt = val
		}
	}
}

// handlePlayerAdded handles new player connection
func (s *WebNowPlayingServer) handlePlayerAdded(playerID int, data string) {
	parsed := parsePlayerData(data)

	player := &Player{
		ID:        playerID,
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}
	applyPlayerData(player, parsed)

	s.playersMu.Lock()
	s.players[playerID] = player
	s.activePlayerID = playerID
	s.playersMu.Unlock()

	log.Printf("Player added: %d (%s) - %s", playerID, player.Name, player.Title)
	s.notifyUpdate(player)
}

// handlePlayerUpdated handles player state update (partial data)
func (s *WebNowPlayingServer) handlePlayerUpdated(playerID int, data string) {
	parsed := parsePlayerData(data)

	s.playersMu.Lock()
	player, ok := s.players[playerID]
	if !ok {
		// Create new player if doesn't exist
		player = &Player{
			ID:        playerID,
			CreatedAt: time.Now().UnixMilli(),
		}
		s.players[playerID] = player
	}
	player.UpdatedAt = time.Now().UnixMilli()
	applyPlayerData(player, parsed)

	// Update active player based on activeAt
	if player.ActiveAt > 0 {
		s.activePlayerID = playerID
	}
	s.playersMu.Unlock()

	s.notifyUpdate(player)
}

// handlePlayerRemoved handles player disconnection
func (s *WebNowPlayingServer) handlePlayerRemoved(playerID int) {
	s.playersMu.Lock()
	delete(s.players, playerID)

	// Find new active player
	if s.activePlayerID == playerID {
		s.activePlayerID = 0
		var maxActiveAt int64
		for id, p := range s.players {
			if p.ActiveAt > maxActiveAt {
				maxActiveAt = p.ActiveAt
				s.activePlayerID = id
			}
		}
	}
	s.playersMu.Unlock()

	log.Printf("Player removed: %d", playerID)

	// Notify with current active player or nil
	s.playersMu.RLock()
	activePlayer := s.players[s.activePlayerID]
	s.playersMu.RUnlock()

	s.notifyUpdate(activePlayer)
}

// notifyUpdate calls the update callback with current active player
func (s *WebNowPlayingServer) notifyUpdate(player *Player) {
	if s.onUpdate != nil {
		s.onUpdate(player.Clone())
	}
}

// SendCommand sends a control command to WebNowPlaying
func (s *WebNowPlayingServer) SendCommand(playerID int, command string, data interface{}) error {
	s.connMu.Lock()
	conn := s.conn
	s.connMu.Unlock()

	if conn == nil {
		return fmt.Errorf("not connected")
	}

	// Format command based on type
	var msg string
	switch command {
	case "STATE":
		msg = fmt.Sprintf("STATE %d %v", playerID, data)
	case "SKIP_NEXT":
		msg = fmt.Sprintf("SKIP_NEXT %d", playerID)
	case "SKIP_PREVIOUS":
		msg = fmt.Sprintf("SKIP_PREVIOUS %d", playerID)
	case "SHUFFLE":
		msg = fmt.Sprintf("SHUFFLE %d %v", playerID, data)
	case "REPEAT":
		msg = fmt.Sprintf("REPEAT %d %v", playerID, data)
	case "RATING":
		msg = fmt.Sprintf("RATING %d %v", playerID, data)
	case "POSITION":
		msg = fmt.Sprintf("POSITION %d %v", playerID, data)
	case "VOLUME":
		msg = fmt.Sprintf("VOLUME %d %v", playerID, data)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	log.Printf("Sending command: %s", msg)
	return conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

// GetActivePlayer returns the current active player
func (s *WebNowPlayingServer) GetActivePlayer() *Player {
	s.playersMu.RLock()
	defer s.playersMu.RUnlock()

	if player, ok := s.players[s.activePlayerID]; ok {
		return player.Clone()
	}
	return nil
}
