package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const autorunRegistryKey = `Software\Microsoft\Windows\CurrentVersion\Run`
const appRegistryName = "RoundSound"

type AutorunManager struct {
	execPath string
}

func NewAutorunManager() (*AutorunManager, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	execPath = filepath.Clean(execPath)

	return &AutorunManager{
		execPath: execPath,
	}, nil
}

func (am *AutorunManager) IsEnabled() (bool, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, autorunRegistryKey, registry.QUERY_VALUE)
	if err != nil {
		return false, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	val, _, err := key.GetStringValue(appRegistryName)
	if err == registry.ErrNotExist {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to read registry value: %w", err)
	}

	return val == am.execPath, nil
}

func (am *AutorunManager) Enable() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, autorunRegistryKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(appRegistryName, am.execPath); err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	log.Printf("[Autorun] Enabled autostart: %s", am.execPath)
	return nil
}

func (am *AutorunManager) Disable() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, autorunRegistryKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.DeleteValue(appRegistryName); err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("failed to delete registry value: %w", err)
	}

	log.Println("[Autorun] Disabled autostart")
	return nil
}
