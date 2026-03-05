package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/jack/pomodoro-timer/internal/timer"
)

func main() {
	p := tea.NewProgram(timer.New())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
