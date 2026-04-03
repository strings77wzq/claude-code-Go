package tui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type streamMsg struct {
	text string
}

type doneMsg struct {
	result string
}

type errorMsg struct {
	err error
}

type message struct {
	role    string
	content string
}

type model struct {
	messages     []message
	input        textinput.Model
	spinner      spinner.Model
	isLoading    bool
	streamBuffer string
	provider     string
	modelName    string
	version      string
	quitting     bool
	agent        AgentInterface
}

type AgentInterface interface {
	Run(ctx context.Context, userInput string, onTextDelta func(text string)) (string, error)
	ClearHistory()
	SetModel(model string)
	Model() string
}

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00ADD8"))
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7ee787")).Bold(true)
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5f56"))
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#8b949e"))
)

func NewModel(agent AgentInterface, version, provider, modelName string) model {
	ti := textinput.New()
	ti.Placeholder = "Type a message or /help for commands..."
	ti.Prompt = promptStyle.Render("go-code> ")
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = 80

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ADD8"))

	return model{
		messages:  make([]message, 0),
		input:     ti,
		spinner:   s,
		agent:     agent,
		provider:  provider,
		modelName: modelName,
		version:   version,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlD:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			return m.handleEnter()
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd

	case streamMsg:
		m.streamBuffer += msg.text
		return m, nil

	case doneMsg:
		if m.streamBuffer != "" {
			m.messages = append(m.messages, message{role: "assistant", content: m.streamBuffer})
		} else if msg.result != "" {
			m.messages = append(m.messages, message{role: "assistant", content: msg.result})
		}
		m.streamBuffer = ""
		m.isLoading = false
		return m, nil

	case errorMsg:
		m.messages = append(m.messages, message{role: "error", content: msg.err.Error()})
		m.streamBuffer = ""
		m.isLoading = false
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	input := m.input.Value()
	m.input.SetValue("")

	if input == "" {
		return m, nil
	}

	if strings.HasPrefix(input, "/") {
		return m.handleCommand(input)
	}

	m.messages = append(m.messages, message{role: "user", content: input})
	m.isLoading = true

	return m, m.runAgent(input)
}

func (m model) runAgent(input string) tea.Cmd {
	streamChan := make(chan streamMsg, 64)
	doneChan := make(chan doneMsg, 1)
	errChan := make(chan errorMsg, 1)

	go func() {
		result, err := m.agent.Run(context.Background(), input, func(text string) {
			streamChan <- streamMsg{text: text}
		})
		if err != nil {
			errChan <- errorMsg{err: err}
		} else {
			doneChan <- doneMsg{result: result}
		}
	}()

	return tea.Batch(
		func() tea.Msg {
			for {
				select {
				case msg := <-streamChan:
					return msg
				case msg := <-doneChan:
					return msg
				case msg := <-errChan:
					return msg
				}
			}
		},
	)
}

func (m model) handleCommand(input string) (tea.Model, tea.Cmd) {
	switch input {
	case "/exit", "/quit":
		m.quitting = true
		return m, tea.Quit
	case "/clear":
		m.messages = nil
		m.agent.ClearHistory()
		return m, nil
	case "/help":
		help := "Commands:\n" +
			"  /help      - Show this help\n" +
			"  /clear     - Clear history\n" +
			"  /model     - Show current model\n" +
			"  /model <n> - Switch model\n" +
			"  /models    - List available models\n" +
			"  /sessions  - List sessions\n" +
			"  /resume <id> - Resume session\n" +
			"  /compact   - Compact context\n" +
			"  /update    - Check for updates\n" +
			"  /exit      - Exit"
		m.messages = append(m.messages, message{role: "system", content: help})
		return m, nil
	case "/model":
		m.messages = append(m.messages, message{role: "system", content: "Current model: " + m.agent.Model()})
		return m, nil
	default:
		if strings.HasPrefix(input, "/model ") {
			newModel := strings.TrimSpace(input[7:])
			m.agent.SetModel(newModel)
			m.modelName = newModel
			m.messages = append(m.messages, message{role: "system", content: "Model switched to: " + newModel})
			return m, nil
		}
		m.messages = append(m.messages, message{role: "system", content: "Unknown command: " + input})
		return m, nil
	}
}

func (m model) View() string {
	if m.quitting {
		return dimStyle.Render("Goodbye!")
	}

	var s string

	s += titleStyle.Render("go-code "+m.version) + "\n"
	s += dimStyle.Render("Model: "+m.modelName+" | Provider: "+m.provider) + "\n"
	s += dimStyle.Render(strings.Repeat("─", 50)) + "\n\n"

	for _, msg := range m.messages {
		switch msg.role {
		case "user":
			s += promptStyle.Render("> ") + msg.content + "\n\n"
		case "assistant":
			s += msg.content + "\n\n"
		case "error":
			s += errorStyle.Render("Error: "+msg.content) + "\n\n"
		case "system":
			s += dimStyle.Render(msg.content) + "\n\n"
		}
	}

	if m.isLoading {
		s += m.spinner.View() + " Thinking...\n\n"
	} else if m.streamBuffer != "" {
		s += m.streamBuffer + "\n\n"
	}

	s += m.input.View()

	return s
}
