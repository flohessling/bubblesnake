package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/flohessling/bubblesnake/internal/model"
)

func main() {
	rand.NewSource(time.Now().UnixNano())
	p := tea.NewProgram(model.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("this did not go well: %v", err)
		os.Exit(1)
	}
}
