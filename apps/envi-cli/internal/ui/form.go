package ui

import (
	"envi/internal/llog"
	"log"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	textInput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TextInput struct {
	questions       []Question
	currentQuestion int
	width           int
	height          int
	answerField     textInput.Model
}

type Question struct {
	Question string
	Answer   string
	EchoMode textInput.EchoMode
}

func NewQuestion(question string) Question {
	return Question{Question: question, Answer: ""}
}
func (q Question) WithEchoMode(echoMode textInput.EchoMode) Question {
	q.EchoMode = echoMode
	return q
}

func (p *TextInput) Init() tea.Cmd {
	return nil
}

func (model *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "ctrl+d":
			return model, tea.Quit
		case "enter":
			model.questions[model.currentQuestion].Answer = model.answerField.Value()
			if model.NextQuestion() {
				model.answerField.SetValue("")
				return model, tea.Quit
			}
			model.answerField.SetValue("")
			return model, nil
		}
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
	}
	model.answerField.EchoMode = model.questions[model.currentQuestion].EchoMode
	model.answerField, cmd = model.answerField.Update(msg)

	return model, cmd
}

func (p *TextInput) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		llog.HelpStyle().Render(p.questions[p.currentQuestion].Question),
		p.answerField.View(),
	)
}

func (p *TextInput) NextQuestion() bool {
	if p.currentQuestion < len(p.questions)-1 {
		p.currentQuestion++
		return false
	}
	p.currentQuestion = 0
	return true
}

func (q Question) Value(question Question) string {
	return question.Answer
}

func NewPrompt(questions []Question) []string {
	answerField := textInput.New()
	answerField.Placeholder = "Type here"
	answerField.Focus()

	p := tea.NewProgram(
		&TextInput{
			questions:   questions,
			answerField: answerField,
		},
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	return F.Pipe1(
		questions,
		A.Map(Question{}.Value),
	)
}
