package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/go-code/internal/agent"
	"github.com/user/go-code/internal/api"
	"github.com/user/go-code/internal/config"
	"github.com/user/go-code/internal/permission"
	"github.com/user/go-code/internal/tool"
	toolinit "github.com/user/go-code/internal/tool/init"
	"github.com/user/go-code/pkg/tty"
)

const version = "0.1.0"

const systemPrompt = "You are an interactive agent that helps users with software engineering tasks. You have access to tools for reading files, editing files, executing shell commands, searching code, and more. Use your tools to complete tasks efficiently and accurately."

func main() {
	// Initialize structured logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Register signal handlers for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received signal, shutting down", "signal", sig.String())
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
	repl := tty.NewREPL(agentInstance, version, cfg.Model)

	// Run REPL - this blocks until exit
	repl.Run()

	logger.Info("REPL exited")
}
