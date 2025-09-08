package main

import (
	"os"
)

func runGUI() {
	installDesktopFile()
	
	app := NewApp()
	os.Exit(app.Run())
}

func shouldRunGUI() bool {
	if len(os.Args) == 1 {
		return true
	}
	
	for _, arg := range os.Args[1:] {
		if arg == "--gui" || arg == "-gui" {
			return true
		}
		if arg == "--cli" || arg == "-cli" {
			return false
		}
	}
	
	return false
}

func filterCLIArgs() []string {
	var filtered []string
	filtered = append(filtered, os.Args[0])
	
	for _, arg := range os.Args[1:] {
		if arg != "--gui" && arg != "-gui" && arg != "--cli" && arg != "-cli" {
			filtered = append(filtered, arg)
		}
	}
	
	return filtered
}
