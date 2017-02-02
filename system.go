package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type System interface {
	Getenv(string) string
	FileExists(string) bool
	DataPath() (string, error)
}

type DefaultSystem struct{}

func (ds DefaultSystem) Getenv(key string) string {
	return os.Getenv(key)
}

func (ds DefaultSystem) FileExists(localPath string) bool {
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (ds DefaultSystem) DataPath() (string, error) {
	var envboxPath string
	if xdgDataHome := ds.Getenv("XDG_DATA_HOME"); len(xdgDataHome) > 0 {
		envboxPath = filepath.Join(xdgDataHome, "envbox")
	} else {
		var home string
		if home = ds.Getenv("HOME"); len(home) == 0 {
			return "", fmt.Errorf("$HOME environment variable not found")
		}
		envboxPath = filepath.Join(home, ".local", "share", "envbox")
	}
	os.MkdirAll(envboxPath, 0755)

	return envboxPath, nil
}
