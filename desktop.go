package main

import (
	"os"
	"path/filepath"
)

func createDesktopFile() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	desktopDir := filepath.Join(homeDir, ".local", "share", "applications")
	os.MkdirAll(desktopDir, 0755)
	
	desktopFile := filepath.Join(desktopDir, "com.shef.app.desktop")
	
	content := `[Desktop Entry]
Name=shef
Comment=Shodan search tool
Exec=shef
Icon=applications-internet
Terminal=false
Type=Application
Categories=Network;Security;
Keywords=shodan;search;security;network;
StartupNotify=true
`
	
	return os.WriteFile(desktopFile, []byte(content), 0644)
}

func installDesktopFile() {
	createDesktopFile()
}
