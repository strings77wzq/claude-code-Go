package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/agent"
	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/config"
	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/skills"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
	toolinit "github.com/strings77wzq/claude-code-Go/internal/tool/init"
	"github.com/strings77wzq/claude-code-Go/pkg/tty"
)

const version = "0.1.0"

const systemPrompt = "You are an interactive agent that helps users with software engineering tasks. You have access to tools for reading files, editing files, executing shell commands, searching code, and more. Use your tools to complete tasks efficiently and accurately."

func main() {
	// Check for --setup flag
	if len(os.Args) > 1 && os.Args[1] == "--setup" {
		if err := SetupWizard(); err != nil {
			fmt.Fprintf(os.Stderr, "Setup failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Initialize structured logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Create context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register signal handlers for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received signal, shutting down", "signal", sig.String())
		cancel()
		// Wait briefly for graceful shutdown
		time.Sleep(100 * time.Millisecond)
		logger.Info("Shutdown complete")
		os.Exit(0)
	}()

	// Step 1: Load configuration
	logger.Info("Loading configuration")
	cfg, err := config.Load(nil)
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	logger.Info("Configuration loaded", "model", cfg.Model, "baseURL", cfg.BaseURL)

	// Step 2: Create API client
	logger.Info("Creating API client")
	client := api.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model)
	logger.Info("API client created")

	// Step 3: Create tool registry
	logger.Info("Creating tool registry")
	registry := tool.NewRegistry()
	logger.Info("Tool registry created")

	// Step 4: Register builtin tools
	logger.Info("Registering builtin tools")
	wd := cfg.WorkingDir
	if wd == "" {
		wd, _ = os.Getwd()
	}
	if err := toolinit.RegisterBuiltinTools(registry, wd); err != nil {
		logger.Error("Failed to register builtin tools", "error", err)
		os.Exit(1)
	}
	logger.Info("Builtin tools registered", "count", len(registry.GetAllDefinitions()))

	// Step 5: Create permission policy
	logger.Info("Creating permission policy")
	policy := permission.NewPolicy(permission.WorkspaceWrite)
	logger.Info("Permission policy created")

	// Step 6: Create agent
	logger.Info("Creating agent")
	agentInstance := agent.NewAgent(client, registry, policy, systemPrompt, cfg.Model)
	logger.Info("Agent started", "model", cfg.Model)

	// Step 7: Create REPL with version and model
	logger.Info("Starting REPL")
	skillsRegistry := skills.NewRegistry()
	if homeDir, err := os.UserHomeDir(); err == nil {
		skillsDir := filepath.Join(homeDir, ".go-code", "skills")
		if loadedSkills, err := skills.LoadSkills(skillsDir); err == nil {
			for _, s := range loadedSkills {
				if err := skillsRegistry.Register(s); err != nil {
					logger.Warn("Failed to register skill", "name", s.Name, "error", err)
				}
			}
		}
	}

	repl := tty.NewREPL(agentInstance, version, cfg.Provider, cfg.Model, skillsRegistry, "~/.go-code/sessions/")

	// Pass external context for graceful shutdown
	repl.SetExternalContext(ctx)

	// Run REPL - this blocks until exit
	repl.Run()

	logger.Info("REPL exited")
}
