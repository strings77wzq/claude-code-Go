package tty

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/strings77wzq/claude-code-Go/internal/agent"
	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/provider/registry"
	"github.com/strings77wzq/claude-code-Go/internal/session"
	"github.com/strings77wzq/claude-code-Go/internal/skills"
	"github.com/strings77wzq/claude-code-Go/internal/update"
)

type AgentInterface interface {
	Run(ctx context.Context, userInput string, onTextDelta func(text string)) (string, error)
}

type REPL struct {
	agent          AgentInterface
	renderer       *Renderer
	version        string
	provider       string
	model          string
	history        []string
	ctx            context.Context
	cancel         context.CancelFunc
	skillsRegistry *skills.Registry
	sessionsDir    string
}

func NewREPL(agent AgentInterface, version string, provider string, model string, skillsRegistry *skills.Registry, sessionsDir string) *REPL {
	ctx, cancel := context.WithCancel(context.Background())
	repl := &REPL{
		agent:          agent,
		renderer:       NewRenderer(),
		version:        version,
		provider:       provider,
		model:          model,
		history:        make([]string, 0, 100),
		ctx:            ctx,
		cancel:         cancel,
		skillsRegistry: skillsRegistry,
		sessionsDir:    sessionsDir,
	}
	return repl
}

func (r *REPL) SetExternalContext(ctx context.Context) {
	r.ctx = ctx
}

func (r *REPL) Run() {
	r.RunWithContext(context.Background())
}

func (r *REPL) RunWithContext(ctx context.Context) {
	r.ctx = ctx
	r.renderer.PrintWelcome(r.version, r.provider, r.model)

	scanner := bufio.NewScanner(os.Stdin)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		for {
			select {
			case <-r.ctx.Done():
				return
			case sig := <-sigChan:
				if sig == syscall.SIGINT {
					r.cancel()
					r.ctx, r.cancel = context.WithCancel(context.Background())
					fmt.Println("\n(Canceled. Press Enter to continue)")
				}
			}
		}
	}()

	for {
		select {
		case <-r.ctx.Done():
			fmt.Println("\nGoodbye!")
			return
		default:
		}

		r.renderer.PrintPrompt()

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		r.addToHistory(input)

		if r.handleSpecialCommand(input) {
			continue
		}

		r.processInput(input)
	}

	fmt.Println("\nGoodbye!")
}

func (r *REPL) handleSpecialCommand(input string) bool {
	if input == "/skills" {
		r.handleSkillsCommand()
		return true
	}

	if input == "/sessions" {
		r.handleSessionsCommand()
		return true
	}

	if remainder, ok := strings.CutPrefix(input, "/resume"); ok {
		sessionID := strings.TrimSpace(remainder)
		if sessionID == "" {
			fmt.Println("Usage: /resume <session-id>")
			return true
		}
		r.handleResumeCommand(sessionID)
		return true
	}

	switch input {
	case "/help":
		r.renderer.PrintHelp()
		return true
	case "/clear":
		// Clear the agent's history if it has a Clear method
		if ah, ok := r.agent.(interface{ GetHistory() *agent.History }); ok {
			ah.GetHistory().Clear()
		}
		fmt.Println("Conversation history cleared")
		return true
	case "/compact":
		if c, ok := r.agent.(interface{ Compact() }); ok {
			c.Compact()
			fmt.Println("Conversation compacted")
		} else {
			fmt.Println("Compact is not available for the current agent")
		}
		return true
	case "/update":
		latestVersion, downloadURL, needsUpdate, err := update.CheckLatest(r.version)
		if err != nil {
			r.renderer.PrintError(err)
			return true
		}

		if !needsUpdate {
			fmt.Printf("Already up to date (%s)\n", r.version)
			return true
		}

		if downloadURL == "" {
			fmt.Printf("Update available (%s) but no download URL found\n", latestVersion)
			return true
		}

		answer, err := r.promptLine(fmt.Sprintf("Update available: %s -> %s. Download and replace binary? [y/N]: ", r.version, latestVersion))
		if err != nil {
			r.renderer.PrintError(err)
			return true
		}

		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Update canceled")
			return true
		}

		binaryPath, err := os.Executable()
		if err != nil {
			r.renderer.PrintError(fmt.Errorf("failed to resolve current binary path: %w", err))
			return true
		}

		if err := update.DownloadAndUpdate(downloadURL, binaryPath); err != nil {
			r.renderer.PrintError(err)
			return true
		}

		fmt.Println("Update successful. Please restart go-code.")
		return true
	case "/exit", "/quit":
		return false
	case "/model":
		if remainder := strings.TrimSpace(strings.TrimPrefix(input, "/model")); remainder != "" {
			if setter, ok := r.agent.(interface{ SetModel(string) }); ok {
				setter.SetModel(remainder)
				r.model = remainder
				fmt.Printf("Model switched to: %s\n", remainder)
			} else {
				fmt.Println("Model switching is not supported for the current agent")
			}
		} else {
			r.renderer.PrintModel(r.model)
		}
		return true
	case "/models":
		r.printAvailableModels()
		return true
	default:
		if r.handleSkillCommand(input) {
			return true
		}
		return false
	}
}

func (r *REPL) processInput(input string) {
	_, err := r.agent.Run(r.ctx, input, func(text string) {
		r.renderer.PrintStreaming(text)
	})
	if err != nil {
		r.renderer.PrintError(err)
	}
	fmt.Println()
}

func (r *REPL) promptLine(prompt string) (string, error) {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", nil
}

func (r *REPL) addToHistory(cmd string) {
	if len(r.history) >= 100 {
		r.history = r.history[1:]
	}
	r.history = append(r.history, cmd)
}

func (r *REPL) handleSkillsCommand() {
	if r.skillsRegistry == nil {
		fmt.Println("Skills registry is not configured")
		return
	}

	skillList := r.skillsRegistry.List()
	if len(skillList) == 0 {
		fmt.Println("No skills available")
		return
	}

	sort.Slice(skillList, func(i, j int) bool {
		return skillList[i].Name < skillList[j].Name
	})

	fmt.Println("Available skills:")
	for _, s := range skillList {
		if s.Description != "" {
			fmt.Printf("  /%s - %s\n", s.Name, s.Description)
			continue
		}
		fmt.Printf("  /%s\n", s.Name)
	}
}

func (r *REPL) handleSkillCommand(input string) bool {
	if !strings.HasPrefix(input, "/") || r.skillsRegistry == nil {
		return false
	}

	skillName := strings.TrimPrefix(input, "/")
	if skillName == "" || strings.Contains(skillName, " ") {
		return false
	}

	prompt, err := r.skillsRegistry.Execute(skillName)
	if err != nil {
		return false
	}

	fmt.Printf("Executing skill /%s\n", skillName)
	r.processInput(prompt)
	return true
}

func (r *REPL) handleSessionsCommand() {
	dir := r.expandSessionsDir()
	sessions, err := session.ListSessions(dir)
	if err != nil {
		r.renderer.PrintError(err)
		return
	}

	if len(sessions) == 0 {
		fmt.Println("No sessions found")
		return
	}

	fmt.Println("Available sessions:")
	for _, s := range sessions {
		fmt.Printf("  %s  model=%s turns=%d started=%s\n", s.ID, s.Model, s.TurnCount, s.StartTime.Format("2006-01-02 15:04:05"))
	}
}

func (r *REPL) handleResumeCommand(sessionID string) {
	dir := r.expandSessionsDir()
	sessions, err := session.ListSessions(dir)
	if err != nil {
		r.renderer.PrintError(err)
		return
	}

	var match *session.SessionInfo
	for i := range sessions {
		if sessions[i].ID == sessionID {
			match = &sessions[i]
			break
		}
	}

	if match == nil {
		fmt.Printf("Session not found: %s\n", sessionID)
		return
	}

	loadedSession, messages, err := session.LoadSession(match.FilePath)
	if err != nil {
		r.renderer.PrintError(err)
		return
	}

	historyAgent, ok := r.agent.(interface{ GetHistory() *agent.History })
	if !ok {
		fmt.Println("Current agent does not support session resume")
		return
	}

	history := historyAgent.GetHistory()
	history.Clear()
	for _, msg := range messages {
		switch msg.Role {
		case "user":
			if err := history.AddUserMessage(msg.Content); err != nil {
				r.renderer.PrintError(err)
				return
			}
		case "assistant":
			if err := history.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: msg.Content}}); err != nil {
				r.renderer.PrintError(err)
				return
			}
		}
	}

	if loadedSession.Model != "" {
		r.model = loadedSession.Model
	}

	fmt.Printf("Resumed session %s with %d messages\n", loadedSession.ID, len(messages))
	if loadedSession.Model != "" {
		fmt.Printf("Session model: %s\n", loadedSession.Model)
	}
}

func (r *REPL) expandSessionsDir() string {
	dir := strings.TrimSpace(r.sessionsDir)
	if dir == "" {
		dir = "~/.go-code/sessions/"
	}

	if strings.HasPrefix(dir, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, strings.TrimPrefix(dir, "~/"))
		}
	}

	return dir
}

func (r *REPL) printAvailableModels() {
	models := registry.GetSupportedModels()

	fmt.Println("Available models:")
	fmt.Println()

	currentProvider := ""
	for _, m := range models {
		if m.Provider != currentProvider {
			currentProvider = m.Provider
			fmt.Printf("  %s:\n", strings.Title(currentProvider))
		}
		fmt.Printf("    %s - %s\n", m.Name, m.Description)
	}

	fmt.Println()
	fmt.Println("Switch model: /model <model-name>")
}
