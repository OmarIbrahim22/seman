package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/OmarIbrahim22/seman/internal/app"
	"github.com/OmarIbrahim22/seman/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	dataDir, err := dataDirectory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error locating data directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating data directory: %v\n", err)
		os.Exit(1)
	}

	store := storage.NewJSONStore(filepath.Join(dataDir, "semester.json"))
	data, found, err := store.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading data: %v\n", err)
		os.Exit(1)
	}
	p := tea.NewProgram(app.NewModel(store, data, found), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// dataDirectory returns ~/.local/share/seman on Linux/macOS-ish systems.
func dataDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "seman"), nil
}
