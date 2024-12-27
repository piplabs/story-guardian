package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	linuxPath  = ".story/geth/guardian"
	darwinPath = "Library/Story/geth/guardian"
)

// GetDefaultPath determines the default file path based on the operating system.
func GetDefaultPath() string {
	userHomeDir, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "linux":
		return filepath.Join(userHomeDir, linuxPath)
	case "darwin":
		return filepath.Join(userHomeDir, darwinPath)
	default:
		log.Fatalf("Unsupported operating system: %s", runtime.GOOS)
		return ""
	}
}
