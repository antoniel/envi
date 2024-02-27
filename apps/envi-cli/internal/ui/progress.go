package ui

import (
	"envi/apps/envi-cli/internal/llog"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time
type MsgComplete struct{}
type MsgQuit struct{}
type model struct {
	progress  progress.Model
	title     string
	didFinish bool
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MsgQuit:
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
			m.didFinish = true
			return m.Update(MsgQuit{})
		}
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.didFinish {
		return "" // Clear the screen
	}
	var boldStyle = llog.HelpStyle().Bold(true)
	return lipgloss.JoinVertical(
		lipgloss.Top,
		llog.HelpStyle().Render(m.title),
		m.progress.View(),
		llog.HelpStyle().Render("Current provider: "+boldStyle.Render("Zipper")),
	)
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
