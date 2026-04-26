// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/config/paths.go

package config

import (
	"errors"
	"os"
	"path/filepath"
)

const DefaultDatabaseName = "orthodox-calendar.db"

func ResolveDBPath(explicitPath string) (string, error) {
	if explicitPath != "" {
		return explicitPath, nil
	}

	if envPath := os.Getenv("ORTHOCAL_DB"); envPath != "" {
		return envPath, nil
	}

	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		return filepath.Join(dataHome, "orthocal", DefaultDatabaseName), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("could not find user home directory")
	}

	return filepath.Join(home, ".local", "share", "orthocal", DefaultDatabaseName), nil
}
