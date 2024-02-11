package llog

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func BgPrimaryColorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#8F30B0")).
		AlignHorizontal(lipgloss.Center)
}

func BgPrimaryColorFullWidth(strs ...string) {
	w, _, _ := term.GetSize(2)
	fmt.Println(
		BgPrimaryColorStyle().
			PaddingLeft(2).
			PaddingRight(2).
			Width(w).
			Render(strs...))
}

func HelpStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
}

func StyleCommand() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AD58B4")).
		Bold(true).
		PaddingLeft(2)
}

func StyleTitle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		MarginBottom(1)
}
