package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time
type MsgComplete struct{}
type model struct {
	progress progress.Model
	title    string
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width
		return m, nil
	case tickMsg:
		cmd := m.progress.IncrPercent(0.01)
		return m, tea.Batch(tickCmd(), cmd)
	case MsgComplete:
		cmd := m.progress.SetPercent(1)
		return m, cmd
	case progress.FrameMsg:
		if m.progress.Percent() == 1 && !m.progress.IsAnimating() {
			fmt.Print("\033[K\033[A\033[K") // Clear the current line and move up
			return m, tea.Quit
		}
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
	return lipgloss.JoinVertical(
		lipgloss.Top,
		helpStyle(m.title),
		m.progress.View())
}

func (t tickMsg) Second() float64 {
	return time.Since(time.Time(t)).Seconds()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func ProgressBar(title string) func() {
	defaultOpts := []progress.Option{
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
	}

	doneCh := make(chan bool)
	m := model{progress: progress.New(defaultOpts...), title: title}
	program := tea.NewProgram(m)
	go goRunProgram(program, doneCh)

	return func() {
		program.Send(MsgComplete{})
		<-doneCh // Just wait until goRunProgram is done
	}
}

func goRunProgram(program *tea.Program, doneCh chan bool) {
	_, err := program.Run()
	if err != nil {
		os.Exit(1)
	}
	doneCh <- true
	close(doneCh)
}
