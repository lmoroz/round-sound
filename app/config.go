package app

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	WindowX int `json:"windowX"`
	WindowY int `json:"windowY"`
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
	cfg := &Config{}
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

	log.Printf("Loaded config: WindowX=%d, WindowY=%d", cfg.WindowX, cfg.WindowY)
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

	log.Printf("Saved config: WindowX=%d, WindowY=%d", c.WindowX, c.WindowY)
	return nil
}
