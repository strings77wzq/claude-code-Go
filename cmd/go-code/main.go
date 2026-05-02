package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/strings77wzq/claude-code-Go/internal/agent"
	"github.com/strings77wzq/claude-code-Go/internal/config"
	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/provider"
	"github.com/strings77wzq/claude-code-Go/internal/provider/registry"
	"github.com/strings77wzq/claude-code-Go/internal/skills"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
	toolinit "github.com/strings77wzq/claude-code-Go/internal/tool/init"
	"github.com/strings77wzq/claude-code-Go/internal/tool/mcp"
	"github.com/strings77wzq/claude-code-Go/pkg/tty"
	"github.com/strings77wzq/claude-code-Go/pkg/tui"
)

const version = "0.3.0"

const systemPrompt = "You are an interactive agent that helps users with software engineering tasks. You have access to tools for reading files, editing files, executing shell commands, searching code, and more. Use your tools to complete tasks efficiently and accurately."

type cliOptions struct {
	legacyRepl   bool
	setupMode    bool
	prompt       string
	outputFormat string
	quiet        bool
	debug        bool
	version      bool
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "doctor" {
		os.Exit(runDoctorCommand(os.Args[2:], os.Stdout, os.Stderr))
	}
	if len(os.Args) > 1 && os.Args[1] == "replay" {
		os.Exit(runReplayCommand(os.Args[2:], os.Stdout, os.Stderr))
	}
	if len(os.Args) > 1 && os.Args[1] == "version" {
		printVersion(os.Stdout)
		return
	}

	flags, opts := newRootFlagSet("go-code", flag.ContinueOnError, os.Stdout)
	if err := flags.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			return
		}
		os.Exit(2)
	}

	if opts.version {
		printVersion(os.Stdout)
		return
	}

	if opts.setupMode {
		if err := SetupWizard(); err != nil {
			fmt.Fprintf(os.Stderr, "Setup failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Initialize structured logging
	logLevel := slog.LevelInfo
	if opts.debug {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
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
		time.Sleep(100 * time.Millisecond)
		logger.Info("Shutdown complete")
		os.Exit(0)
	}()

	// Load configuration
	logger.Info("Loading configuration")
	cfg, err := config.Load(nil)
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	logger.Info("Configuration loaded", "model", cfg.Model, "baseURL", cfg.BaseURL)

	// Create API client via provider registry
	logger.Info("Creating API client")
	resolvedProvider, err := registry.ResolveConfig(cfg.Provider, cfg.BaseURL, cfg.Model, cfg.APIKey)
	if err != nil {
		logger.Error("Invalid provider configuration", "error", err)
		os.Exit(1)
	}
	cfg.Provider = resolvedProvider.Provider
	cfg.BaseURL = resolvedProvider.BaseURL
	cfg.Model = resolvedProvider.Model
	apiClient := registry.SelectProviderFor(cfg.Provider, cfg.APIKey, cfg.BaseURL, cfg.Model)
	client := provider.NewApiClientAdapter(apiClient, func(model string) provider.Provider {
		targetProvider := registry.DetectProvider(model)
		baseURL := cfg.BaseURL
		if targetProvider != cfg.Provider {
			baseURL = registry.DefaultBaseURL(targetProvider)
		}
		return registry.SelectProviderFor(targetProvider, cfg.APIKey, baseURL, model)
	})
	logger.Info("API client created", "provider", apiClient.Name())

	// Create tool registry
	logger.Info("Creating tool registry")
	toolRegistry := tool.NewRegistry()
	logger.Info("Tool registry created")

	// Register builtin tools
	logger.Info("Registering builtin tools")
	wd := cfg.WorkingDir
	if wd == "" {
		wd, _ = os.Getwd()
	}
	if err := toolinit.RegisterBuiltinTools(toolRegistry, wd); err != nil {
		logger.Error("Failed to register builtin tools", "error", err)
		os.Exit(1)
	}
	logger.Info("Builtin tools registered", "count", len(toolRegistry.GetAllDefinitions()))

	// Initialize MCP servers and register their tools
	logger.Info("Initializing MCP servers")
	mcpManager := mcp.NewMcpManager()
	mcpConfigPath := mcp.GetDefaultMcpConfigPath()
	if mcpConfigPath != "" {
		if mcpConfigs, err := mcp.LoadMcpConfigs(mcpConfigPath); err == nil {
			mcpManager.InitializeAndRegister(mcpConfigs, toolRegistry)
		} else {
			logger.Debug("No MCP config found, skipping", "path", mcpConfigPath)
		}
	}
	defer mcpManager.Close()
	logger.Info("MCP initialization complete", "tools", len(toolRegistry.GetAllDefinitions()))

	// Create permission policy
	logger.Info("Creating permission policy")
	policy := permission.NewPolicy(permission.WorkspaceWrite)
	logger.Info("Permission policy created")

	// Create agent
	logger.Info("Creating agent")
	agentInstance := agent.NewAgent(client, toolRegistry, policy, systemPrompt, cfg.Model)
	agentInstance.SetPermissionPrompter(permission.NewStdinPrompter(bufio.NewReader(os.Stdin), os.Stdout))
	logger.Info("Agent started", "model", cfg.Model)

	// Non-interactive mode: run single prompt and exit
	if opts.prompt != "" {
		result, err := agentInstance.Run(ctx, opts.prompt, func(text string) {
			if !opts.quiet {
				fmt.Print(text)
			}
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		switch opts.outputFormat {
		case "json":
			output := map[string]string{"result": result}
			data, _ := json.Marshal(output)
			fmt.Println(string(data))
		default:
			if !opts.quiet {
				fmt.Println()
			}
		}
		return
	}

	// Load skills with validation warnings
	skillsRegistry := skills.NewRegistry()
	if homeDir, err := os.UserHomeDir(); err == nil {
		skillsDir := filepath.Join(homeDir, ".go-code", "skills")
		if result, err := skills.LoadSkillsWithWarnings(skillsDir); err == nil {
			for _, w := range result.Warnings {
				logger.Warn("Invalid skill file", "file", w.File, "reason", w.Reason)
			}
			for _, s := range result.Skills {
				if err := skillsRegistry.Register(s); err != nil {
					logger.Warn("Failed to register skill", "name", s.Name, "error", err)
				}
			}
		}
	}

	// Use legacy REPL or new bubbletea TUI
	if opts.legacyRepl {
		logger.Info("Starting legacy REPL")
		repl := tty.NewREPL(agentInstance, version, cfg.Provider, cfg.Model, skillsRegistry, "~/.claude-code-go/sessions/")
		repl.SetExternalContext(ctx)
		repl.Run()
	} else {
		logger.Info("Starting bubbletea TUI")
		tuiModel := tui.NewModel(agentInstance, version, cfg.Provider, cfg.Model, opts.debug)
		p := tea.NewProgram(tuiModel, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			logger.Error("TUI error", "error", err)
			os.Exit(1)
		}
	}

	logger.Info("REPL exited")
}

func newRootFlagSet(name string, errorHandling flag.ErrorHandling, output io.Writer) (*flag.FlagSet, *cliOptions) {
	opts := &cliOptions{}
	flags := flag.NewFlagSet(name, errorHandling)
	flags.SetOutput(output)
	flags.BoolVar(&opts.legacyRepl, "legacy-repl", false, "Use the old bufio-based REPL instead of the default interactive TUI")
	flags.BoolVar(&opts.setupMode, "setup", false, "Run setup wizard")
	flags.StringVar(&opts.prompt, "p", "", "Run a single prompt and exit (non-interactive mode)")
	flags.StringVar(&opts.outputFormat, "f", "text", "Output format for non-interactive mode (text, json)")
	flags.BoolVar(&opts.quiet, "q", false, "Hide spinner in non-interactive mode")
	flags.BoolVar(&opts.debug, "debug", false, "Enable debug logging to stderr")
	flags.BoolVar(&opts.version, "version", false, "Print version and exit")
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "Usage: %s [options]\n", name)
		fmt.Fprintf(flags.Output(), "       %s doctor [--offline]\n", name)
		fmt.Fprintf(flags.Output(), "       %s replay [--evidence] [latest|session-id|file]\n", name)
		fmt.Fprintf(flags.Output(), "       %s version\n\n", name)
		fmt.Fprintln(flags.Output(), "Entrypoints:")
		fmt.Fprintln(flags.Output(), "  interactive mode    default TUI when no prompt is provided")
		fmt.Fprintln(flags.Output(), "  setup               configure provider credentials")
		fmt.Fprintln(flags.Output(), "  doctor              validate local runtime readiness")
		fmt.Fprintln(flags.Output(), "  prompt mode         run one prompt with -p")
		fmt.Fprintln(flags.Output(), "  JSON output         use -f json with prompt mode")
		fmt.Fprintln(flags.Output(), "  quiet mode          suppress streaming text with -q")
		fmt.Fprintln(flags.Output(), "  debug mode          emit debug logs with -debug")
		fmt.Fprintln(flags.Output(), "  permission mode     default WorkspaceWrite; non-interactive prompts fail closed")
		fmt.Fprintln(flags.Output(), "  version             print go-code version")
		fmt.Fprintln(flags.Output(), "\nOptions:")
		flags.PrintDefaults()
	}
	return flags, opts
}

func printVersion(w io.Writer) {
	fmt.Fprintf(w, "go-code %s\n", version)
}
