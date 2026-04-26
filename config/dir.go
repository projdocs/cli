package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetConfigDirPath() string {
	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "projdocs")
		}
		return filepath.Join(os.Getenv("USERPROFILE"), "projdocs")
	default:
		return "/opt/projdocs"
	}
}

func GetConfigDir() (string, error) {
	dir := GetConfigDirPath()
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return "", fmt.Errorf("could not create config dir %q: %w", dir, err)
			}
			return dir, nil
		}
		return "", fmt.Errorf("could not stat config dir %q: %w", dir, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("config path %q exists but is not a directory", dir)
	}
	return dir, nil
}

func MustGetConfigDir() string {
	dir, err := GetConfigDir()
	if err != nil {
		panic(err)
	}
	return dir
}
