package command

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/agent"
	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/provider/registry"
	"github.com/strings77wzq/claude-code-Go/internal/session"
	"github.com/strings77wzq/claude-code-Go/internal/update"
)

type Handler struct {
	Agent       any
	Version     string
	Model       string
	SessionsDir string
	CheckUpdate func(currentVersion string) (latestVersion string, downloadURL string, needsUpdate bool, err error)
}

type Result struct {
	Handled bool
	Quit    bool
	Message string
	Model   string
}

func (h Handler) Handle(input string) Result {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, "/") {
		return Result{}
	}

	switch {
	case input == "/help":
		return h.message(helpText())
	case input == "/clear":
		if clearer, ok := h.Agent.(interface{ ClearHistory() }); ok {
			clearer.ClearHistory()
			return h.message("Conversation history cleared")
		}
		return h.message("Conversation history cannot be cleared by the current agent")
	case input == "/exit" || input == "/quit":
		return Result{Handled: true, Quit: true}
	case input == "/model":
		model := h.Model
		if getter, ok := h.Agent.(interface{ Model() string }); ok {
			model = getter.Model()
		}
		if strings.TrimSpace(model) == "" {
			model = "unknown"
		}
		return h.message("Current model: " + model)
	case strings.HasPrefix(input, "/model "):
		return h.switchModel(strings.TrimSpace(strings.TrimPrefix(input, "/model ")))
	case input == "/models":
		return h.message(formatModels())
	case input == "/sessions":
		return h.listSessions()
	case strings.HasPrefix(input, "/resume"):
		sessionID := strings.TrimSpace(strings.TrimPrefix(input, "/resume"))
		if sessionID == "" {
			return h.message("Usage: /resume <session-id>")
		}
		return h.resumeSession(sessionID)
	case input == "/compact":
		if compactor, ok := h.Agent.(interface{ Compact() }); ok {
			compactor.Compact()
			return h.message("Conversation compacted")
		}
		return h.message("Compact is not available for the current agent")
	case input == "/update":
		return h.checkForUpdate()
	case input == "/permissions":
		return h.message("Permission mode details are not exposed yet. Safe-default approval flow is tracked in PARITY.md.")
	default:
		return h.message("Unknown command: " + input + "\nType /help for available commands.")
	}
}

func (h Handler) message(text string) Result {
	return Result{Handled: true, Message: text, Model: h.Model}
}

func (h Handler) switchModel(model string) Result {
	if model == "" {
		return h.message("Usage: /model <model-name>")
	}

	setter, ok := h.Agent.(interface{ SetModel(string) })
	if !ok {
		return h.message("Model switching is not supported for the current agent")
	}

	if !isKnownModel(model) {
		return h.message(fmt.Sprintf("Unsupported model: %s\nUse /models to list supported models. Keeping current model: %s", model, h.currentModel()))
	}

	setter.SetModel(model)
	return Result{
		Handled: true,
		Message: "Model switched to: " + model,
		Model:   model,
	}
}

func (h Handler) currentModel() string {
	model := h.Model
	if getter, ok := h.Agent.(interface{ Model() string }); ok {
		model = getter.Model()
	}
	if strings.TrimSpace(model) == "" {
		return "unknown"
	}
	return model
}

func isKnownModel(model string) bool {
	return registry.IsKnownModel(model)
}

func formatModels() string {
	models := registry.GetSupportedModels()
	var b strings.Builder
	b.WriteString("Available models:\n")

	currentProvider := ""
	for _, model := range models {
		if model.Deprecated {
			continue
		}
		if model.Provider != currentProvider {
			currentProvider = model.Provider
			b.WriteString("\n")
			b.WriteString("  ")
			b.WriteString(strings.Title(currentProvider))
			b.WriteString(":\n")
		}
		b.WriteString("    ")
		b.WriteString(model.Name)
		b.WriteString(" - ")
		b.WriteString(model.Description)
		b.WriteString("\n")
	}
	b.WriteString("\nSwitch model: /model <model-name>")
	return b.String()
}

func (h Handler) listSessions() Result {
	dir := expandSessionsDir(h.SessionsDir)
	sessions, err := session.ListSessions(dir)
	if err != nil {
		if os.IsNotExist(err) || strings.Contains(err.Error(), "does not exist") {
			return h.message("No sessions found in " + dir)
		}
		return h.message("Failed to list sessions: " + err.Error())
	}
	if len(sessions) == 0 {
		return h.message("No sessions found in " + dir)
	}

	var b strings.Builder
	b.WriteString("Available sessions:\n")
	for _, s := range sessions {
		b.WriteString(fmt.Sprintf("  %s  model=%s turns=%d started=%s\n",
			s.ID,
			s.Model,
			s.TurnCount,
			s.StartTime.Format("2006-01-02 15:04:05"),
		))
	}
	return h.message(strings.TrimRight(b.String(), "\n"))
}

func (h Handler) resumeSession(sessionID string) Result {
	dir := expandSessionsDir(h.SessionsDir)
	sessions, err := session.ListSessions(dir)
	if err != nil {
		return h.message("Failed to list sessions: " + err.Error())
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].StartTime.After(sessions[j].StartTime)
	})

	var match *session.SessionInfo
	for i := range sessions {
		if sessions[i].ID == sessionID {
			match = &sessions[i]
			break
		}
	}
	if match == nil {
		return h.message("Session not found: " + sessionID)
	}

	loadedSession, messages, err := session.LoadSession(match.FilePath)
	if err != nil {
		return h.message("Failed to load session: " + err.Error())
	}

	historyAgent, ok := h.Agent.(interface{ GetHistory() *agent.History })
	if !ok {
		return h.message("Current agent does not support session resume")
	}

	history := historyAgent.GetHistory()
	history.Clear()
	for _, msg := range messages {
		switch msg.Role {
		case "user":
			if err := history.AddUserMessage(msg.Content); err != nil {
				return h.message("Failed to restore user message: " + err.Error())
			}
		case "assistant":
			if err := history.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: msg.Content}}); err != nil {
				return h.message("Failed to restore assistant message: " + err.Error())
			}
		}
	}

	message := fmt.Sprintf("Resumed session %s with %d messages", loadedSession.ID, len(messages))
	if loadedSession.Model != "" {
		message += "\nSession model: " + loadedSession.Model
	}
	return Result{Handled: true, Message: message, Model: loadedSession.Model}
}

func (h Handler) checkForUpdate() Result {
	check := h.CheckUpdate
	if check == nil {
		check = update.CheckLatest
	}

	latestVersion, downloadURL, needsUpdate, err := check(h.Version)
	if err != nil {
		return h.message("Update check failed: " + err.Error())
	}
	if !needsUpdate {
		return h.message("Already up to date (" + h.Version + ")")
	}
	if downloadURL == "" {
		return h.message("Update available (" + latestVersion + ") but no download URL was found")
	}
	return h.message(fmt.Sprintf("Update available: %s -> %s\nDownload: %s", h.Version, latestVersion, downloadURL))
}

func expandSessionsDir(dir string) string {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		dir = filepath.Join("~", ".claude-code-go", "sessions")
	}
	if strings.HasPrefix(dir, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(dir, "~/"))
		}
	}
	return dir
}

func helpText() string {
	return "Commands:\n" +
		"  /help        - Show this help\n" +
		"  /clear       - Clear conversation history\n" +
		"  /model       - Show current model\n" +
		"  /model <n>   - Switch model\n" +
		"  /models      - List available models\n" +
		"  /sessions    - List sessions\n" +
		"  /resume <id> - Resume session\n" +
		"  /compact     - Compact context\n" +
		"  /permissions - Show permission status\n" +
		"  /update      - Check for updates\n" +
		"  /exit        - Exit"
}
