package tty

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/user/go-code/internal/agent"
)

type AgentInterface interface {
	Run(ctx context.Context, userInput string, onTextDelta func(text string)) (string, error)
}

type REPL struct {
	agent    AgentInterface
	renderer *Renderer
	version  string
	model    string
	history  []string
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewREPL(agent AgentInterface, version string, model string) *REPL {
	ctx, cancel := context.WithCancel(context.Background())
	return &REPL{
		agent:    agent,
		renderer: NewRenderer(),
		version:  version,
		model:    model,
		history:  make([]string, 0, 100),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (r *REPL) Run() {
	r.renderer.PrintWelcome(r.version)

	scanner := bufio.NewScanner(os.Stdin)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		for {
			sig := <-sigChan
			if sig == syscall.SIGINT {
				r.cancel()
				r.ctx, r.cancel = context.WithCancel(context.Background())
				fmt.Println("\n(Canceled. Press Enter to continue)")
			}
		}
	}()

	for {
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

		r.processInput(input, scanner)
	}

	fmt.Println("\nGoodbye!")
}

func (r *REPL) readLine(prompt string) (string, error) {
	r.renderer.PrintPrompt()
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", nil
}

func (r *REPL) handleSpecialCommand(input string) bool {
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
	case "/exit", "/quit":
		return false
	case "/model":
		r.renderer.PrintModel(r.model)
		return true
	default:
		return false
	}
}

func (r *REPL) processInput(input string, scanner *bufio.Scanner) {
	_, err := r.agent.Run(r.ctx, input, func(text string) {
		r.renderer.PrintStreaming(text)
	})
	if err != nil {
		r.renderer.PrintError(err)
	}
	fmt.Println()
}

func (r *REPL) addToHistory(cmd string) {
	if len(r.history) >= 100 {
		r.history = r.history[1:]
	}
	r.history = append(r.history, cmd)
}
