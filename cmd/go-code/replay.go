package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/strings77wzq/claude-code-Go/internal/session"
)

func runReplayCommand(args []string, stdout, stderr io.Writer) int {
	target := "latest"
	if len(args) > 0 && args[0] != "" {
		target = args[0]
	}

	filePath, err := resolveReplayTarget(target)
	if err != nil {
		fmt.Fprintf(stderr, "Replay failed: %v\n", err)
		return 1
	}

	events, err := session.ReplaySessionFile(filePath)
	if err != nil {
		fmt.Fprintf(stderr, "Replay failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Replay: %s\n", filePath)
	if formatted := session.FormatReplay(events); formatted != "" {
		fmt.Fprintln(stdout, formatted)
	}
	return 0
}

func resolveReplayTarget(target string) (string, error) {
	if target == "latest" || target == "" {
		return session.GetLastSessionFilePath(defaultSessionsDir())
	}
	if info, err := os.Stat(target); err == nil && !info.IsDir() {
		return target, nil
	}

	sessions, err := session.ListSessions(defaultSessionsDir())
	if err != nil {
		return "", err
	}
	for _, s := range sessions {
		if s.ID == target {
			return s.FilePath, nil
		}
	}
	return "", fmt.Errorf("session not found: %s", target)
}

func defaultSessionsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".claude-code-go", "sessions")
	}
	return filepath.Join(home, ".claude-code-go", "sessions")
}
