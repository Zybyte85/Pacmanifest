package config

import (
	"os"
	"path/filepath"
)

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mypm")
}
