package app

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// DefaultWNPPort is the default WebNowPlaying port (same as Rainmeter adapter)
const DefaultWNPPort = 8974

// Config holds application configuration
type Config struct {
	WindowX int `json:"windowX"`
	WindowY int `json:"windowY"`
	WNPPort int `json:"wnpPort"`
}

// getConfigPath returns the path to config file
func getConfigPath() string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		appData = "."
	}
	configDir := filepath.Join(appData, "round-sound")
	os.MkdirAll(configDir, 0755)
	return filepath.Join(configDir, "config.json")
}

// LoadConfig loads configuration from file
func LoadConfig() *Config {
	cfg := &Config{
		WNPPort: DefaultWNPPort, // Default port
	}
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("Config not found, using defaults: %v", err)
		return cfg
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		log.Printf("Failed to parse config: %v", err)
		return &Config{}
	}

	// Ensure WNPPort has a valid default
	if cfg.WNPPort == 0 {
		cfg.WNPPort = DefaultWNPPort
	}

	log.Printf("Loaded config: WindowX=%d, WindowY=%d, WNPPort=%d", cfg.WindowX, cfg.WindowY, cfg.WNPPort)
	return cfg
}

// Save saves configuration to file
func (c *Config) Save() error {
	configPath := getConfigPath()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return err
	}

	log.Printf("Saved config: WindowX=%d, WindowY=%d, WNPPort=%d", c.WindowX, c.WindowY, c.WNPPort)
	return nil
}
