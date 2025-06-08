// Package main provides the GUI entry point for the Mark Master Sheet Consolidator.
// This creates a user-friendly graphical interface using the Fyne framework.
package main

import (
	"mark-master-sheet/internal/gui"
)

func main() {
	app := gui.NewApp()
	app.Run()
}
