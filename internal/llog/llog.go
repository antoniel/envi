package llog

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type tokens struct {
	PrimaryColor           string
	ForegroundColor        string
	BackgroundColor        string
	SuccessColor           string
	ErrorColor             string
	HintColor              string
	CommandForegroundColor string
}

var Tokens = tokens{
	PrimaryColor:           "#8F30B0",
	ForegroundColor:        "#FAFAFA",
	BackgroundColor:        "#7D56F4",
	SuccessColor:           "#C1F3AB",
	HintColor:              "#626262",
	CommandForegroundColor: "#AD58B4",
	ErrorColor:             "#F15C93",
}

func BgPrimaryColorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(Tokens.ForegroundColor)).
		Background(lipgloss.Color(Tokens.BackgroundColor)).
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
	return lipgloss.NewStyle().Foreground(lipgloss.Color(Tokens.HintColor))
}

func StyleCommand() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(Tokens.CommandForegroundColor)).
		Bold(true)
}

func StyleTitle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(Tokens.ForegroundColor)).
		Background(lipgloss.Color(Tokens.BackgroundColor)).
		Padding(0, 1).
		MarginBottom(1)
}
