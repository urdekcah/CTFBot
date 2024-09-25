package main

import (
	"fmt"
	"os"

	"ctfbot.urdekcah.ru/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{ReportTimestamp: true})
)

func main() {

	p := tea.NewProgram(
		ui.NewModel(logger),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
