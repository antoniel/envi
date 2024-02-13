package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	progress progress.Model
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
	return m.progress.View()
}

type tickMsg time.Time

func (t tickMsg) Second() float64 {
	return time.Since(time.Time(t)).Seconds()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type ProgressBarCmd struct {
	program *tea.Program
}

type MsgComplete struct{}

func (p *ProgressBarCmd) Complete() {
	p.program.Send(MsgComplete{})
}

func ProgressBarProgram() *tea.Program {
	defaultOpts := []progress.Option{
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
	}
	m := model{
		progress: progress.New(defaultOpts...),
	}
	program := tea.NewProgram(m)
	return program
}
