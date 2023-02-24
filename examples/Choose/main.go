package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/lipgloss"
	"github.com/yardbirdsax/bubblewrap"
)

func main() {
	options := []string{
		"red",
		"green",
		"blue",
	}
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3170a9"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("fff000"))
	selectedItemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#43aa3d"))
	choice, err := bubblewrap.Choose(options, bubblewrap.WithCursorStyle(cursorStyle), bubblewrap.WithItemStyle(itemStyle), bubblewrap.WithSelectedItemStyle(selectedItemStyle))
	if err != nil {
		switch err.(type) {
		case bubblewrap.CancelError:
			log.Fatal("[warn] user canceled")
		default:
			log.Fatal(fmt.Errorf("error getting choice: %w", err))
		}
	}
	log.Printf("your choices are: %s", choice)
}
