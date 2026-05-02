package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/strings77wzq/claude-code-Go/internal/command"
)

type streamMsg struct {
	requestID string
	text      string
}

type doneMsg struct {
	requestID string
	result    string
}

type errorMsg struct {
	requestID string
	err       error
}

type message struct {
	role    string
	content string
}

type connectionStatusMsg struct {
	requestID  string
	text       string
	elapsedStr string
}

type agentRunEvents struct {
	requestID string
	stream    <-chan streamMsg
	done      <-chan doneMsg
	errs      <-chan errorMsg
	startTime time.Time
	ticker    *time.Ticker
}

type model struct {
	messages        []message
	input           textinput.Model
	spinner         spinner.Model
	isLoading       bool
	streamBuffer    string
	connectionMsg   string
	elapsedTime     string
	provider        string
	modelName       string
	version         string
	quitting        bool
	agent           AgentInterface
	debug           bool
	latency         time.Duration
	tokenUsage      TokenUsage
	activeCancel    context.CancelFunc
	activeRequestID string
	activeEvents    *agentRunEvents
	nextRequestSeq  int
}

type TokenUsage struct {
	InputTokens  int
	OutputTokens int
}

type AgentInterface interface {
	Run(ctx context.Context, userInput string, onTextDelta func(text string)) (string, error)
	ClearHistory()
	SetModel(model string)
	Model() string
}

var (
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00ADD8"))
	promptStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#7ee787")).Bold(true)
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5f56"))
	dimStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#8b949e"))
	assistantStyle = lipgloss.NewStyle().PaddingLeft(0)
	separatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#30363d"))
)

func NewModel(agent AgentInterface, version, provider, modelName string, debug bool) model {
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
		debug:     debug,
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
			if m.activeCancel != nil {
				m.activeCancel()
				m.activeCancel = nil
				m.activeEvents = nil
				m.activeRequestID = ""
				m.isLoading = false
				m.streamBuffer = ""
			}
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			return m.handleEnter()
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd

	case streamMsg:
		if msg.requestID != "" && msg.requestID != m.activeRequestID {
			return m, nil
		}
		m.streamBuffer += msg.text
		m.connectionMsg = ""
		return m, m.waitAgent()

	case connectionStatusMsg:
		if msg.requestID != "" && msg.requestID != m.activeRequestID {
			return m, nil
		}
		m.connectionMsg = msg.text
		m.elapsedTime = msg.elapsedStr
		return m, m.waitAgent()

	case doneMsg:
		if msg.requestID != "" && msg.requestID != m.activeRequestID {
			return m, nil
		}
		if m.streamBuffer != "" {
			m.messages = append(m.messages, message{role: "assistant", content: m.streamBuffer})
		} else if msg.result != "" {
			m.messages = append(m.messages, message{role: "assistant", content: msg.result})
		}
		m.streamBuffer = ""
		m.isLoading = false
		m.elapsedTime = ""
		m.finishActiveRequest()
		return m, nil

	case errorMsg:
		if msg.requestID != "" && msg.requestID != m.activeRequestID {
			return m, nil
		}
		errStr := msg.err.Error()
		if msg.err == context.DeadlineExceeded {
			errStr = "Request timed out after 5 minutes. Check your network connection and API key."
		}
		m.messages = append(m.messages, message{role: "error", content: errStr})
		m.streamBuffer = ""
		m.isLoading = false
		m.finishActiveRequest()
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
	m.nextRequestSeq++
	requestID := fmt.Sprintf("req-%d", m.nextRequestSeq)
	m.activeRequestID = requestID
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	m.activeCancel = cancel

	events := m.startAgentRun(input, requestID, ctx, cancel)
	m.activeEvents = events
	return m, tea.Batch(m.spinner.Tick, m.waitAgent())
}

func (m model) runAgent(input string) tea.Cmd {
	requestID := "req-direct"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	events := m.startAgentRun(input, requestID, ctx, cancel)
	m.activeEvents = events
	return tea.Batch(m.spinner.Tick, m.waitAgent())
}

func (m model) startAgentRun(input string, requestID string, ctx context.Context, cancel context.CancelFunc) *agentRunEvents {
	streamChan := make(chan streamMsg, 64)
	doneChan := make(chan doneMsg, 1)
	errChan := make(chan errorMsg, 1)

	go func() {
		defer cancel()
		result, err := m.agent.Run(ctx, input, func(text string) {
			select {
			case streamChan <- streamMsg{requestID: requestID, text: text}:
			case <-ctx.Done():
			}
		})
		if err != nil {
			errChan <- errorMsg{requestID: requestID, err: err}
		} else {
			doneChan <- doneMsg{requestID: requestID, result: result}
		}
	}()

	return &agentRunEvents{
		requestID: requestID,
		stream:    streamChan,
		done:      doneChan,
		errs:      errChan,
		startTime: time.Now(),
		ticker:    time.NewTicker(500 * time.Millisecond),
	}
}

func (m model) waitAgent() tea.Cmd {
	events := m.activeEvents
	if events == nil {
		return nil
	}
	return func() tea.Msg {
		for {
			select {
			case msg := <-events.stream:
				return msg
			case msg := <-events.done:
				events.ticker.Stop()
				return msg
			case msg := <-events.errs:
				events.ticker.Stop()
				return msg
			case <-events.ticker.C:
				elapsed := time.Since(events.startTime)
				elapsedStr := elapsed.Round(time.Second).String()
				if elapsed > 5*time.Minute {
					events.ticker.Stop()
					return errorMsg{requestID: events.requestID, err: context.DeadlineExceeded}
				} else if elapsed > 30*time.Second {
					return connectionStatusMsg{requestID: events.requestID, text: "Still connecting... check your network or API key", elapsedStr: elapsedStr}
				} else if elapsed > 3*time.Second {
					return connectionStatusMsg{requestID: events.requestID, text: "Connecting to API...", elapsedStr: elapsedStr}
				} else if elapsed > 500*time.Millisecond {
					return connectionStatusMsg{requestID: events.requestID, text: "", elapsedStr: elapsedStr}
				}
			}
		}
	}
}

func (m *model) finishActiveRequest() {
	if m.activeEvents != nil && m.activeEvents.ticker != nil {
		m.activeEvents.ticker.Stop()
	}
	m.activeCancel = nil
	m.activeEvents = nil
	m.activeRequestID = ""
	m.elapsedTime = ""
	m.connectionMsg = ""
}

func (m model) handleCommand(input string) (tea.Model, tea.Cmd) {
	result := command.Handler{
		Agent:       m.agent,
		Version:     m.version,
		Model:       m.modelName,
		SessionsDir: "~/.claude-code-go/sessions/",
	}.Handle(input)

	if !result.Handled {
		m.messages = append(m.messages, message{role: "system", content: "Unknown command: " + input})
		return m, nil
	}
	if result.Model != "" {
		m.modelName = result.Model
	}
	if result.Quit {
		m.quitting = true
		return m, tea.Quit
	}
	if input == "/clear" {
		m.messages = nil
	}
	if result.Message != "" {
		m.messages = append(m.messages, message{role: "system", content: result.Message})
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return dimStyle.Render("Goodbye!")
	}

	var s string

	s += titleStyle.Render("go-code "+m.version) + "\n"
	s += dimStyle.Render("Model: "+m.modelName+" | Provider: "+m.provider) + "\n"
	s += separatorStyle.Render(strings.Repeat("─", 50)) + "\n\n"

	for _, msg := range m.messages {
		switch msg.role {
		case "user":
			s += promptStyle.Render("❯ ") + msg.content + "\n\n"
		case "assistant":
			s += assistantStyle.Render(msg.content) + "\n\n"
		case "error":
			s += errorStyle.Render("✗ Error: "+msg.content) + "\n\n"
		case "system":
			s += dimStyle.Render(msg.content) + "\n\n"
		}
	}

	if m.isLoading {
		loadingText := "Thinking..."
		if m.elapsedTime != "" {
			loadingText += " (" + m.elapsedTime + ")"
		}
		s += m.spinner.View() + " " + loadingText + "\n"
		if m.connectionMsg != "" {
			s += dimStyle.Render("  "+m.connectionMsg) + "\n"
		}
		s += "\n"
	} else if m.streamBuffer != "" {
		s += m.streamBuffer + "\n\n"
	}

	s += m.input.View()

	// Debug status bar
	if m.debug {
		s += "\n" + separatorStyle.Render(strings.Repeat("─", 50)) + "\n"
		debugInfo := "DEBUG | Model: " + m.modelName
		if m.isLoading && m.latency > 0 {
			debugInfo += " | Latency: " + m.latency.Round(time.Millisecond).String()
		}
		if m.tokenUsage.InputTokens > 0 || m.tokenUsage.OutputTokens > 0 {
			debugInfo += fmt.Sprintf(" | Tokens: in=%d out=%d", m.tokenUsage.InputTokens, m.tokenUsage.OutputTokens)
		}
		s += dimStyle.Render(debugInfo)
	}

	return s
}
