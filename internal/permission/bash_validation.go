package permission

import (
	"errors"
	"strings"
)

// CommandCategory represents the category of a bash command
type CommandCategory int

const (
	// CmdUnknown represents an unrecognized command
	CmdUnknown CommandCategory = iota
	// CmdReadOnly represents a safe, read-only command
	CmdReadOnly
	// CmdWrite represents a command that modifies files but is generally safe
	CmdWrite
	// CmdDangerous represents a potentially harmful command
	CmdDangerous
)

// String returns the string representation of CommandCategory
func (c CommandCategory) String() string {
	switch c {
	case CmdReadOnly:
		return "ReadOnly"
	case CmdWrite:
		return "Write"
	case CmdDangerous:
		return "Dangerous"
	default:
		return "Unknown"
	}
}

// ReadOnlyCommands is a whitelist of read-only commands
var ReadOnlyCommands = []string{
	"ls", "cat", "grep", "find", "wc", "head", "tail", "echo", "pwd", "tree",
	"stat", "file", "du", "diff", "sort", "uniq", "cut", "tr", "tee", "man",
	"which", "whereis",
}

// DangerousPatterns is a list of dangerous command patterns
var DangerousPatterns = []string{
	"rm -rf /",
	"rm -rf ~",
	"curl | bash",
	"wget | bash",
	"sudo",
	"dd if=",
	"mkfs",
	"fdisk",
	"chmod 777 /",
	"chown root",
	"> /dev/sda",
	":(){:|:&};:",
	"fork bomb",
}

// WriteCommands is a list of commands that modify files but are generally safe in workspace
var WriteCommands = []string{
	"mkdir", "touch", "cp", "mv", "rm", "sed", "awk", "chmod", "chown",
}

// errDangerousCommand is returned when a dangerous command is detected
var errDangerousCommand = errors.New("dangerous command detected")

// errPathInjection is returned when a path injection attempt is detected
var errPathInjection = errors.New("path injection detected")

// errInvalidPath is returned when a path is invalid or escapes workspace
var errInvalidPath = errors.New("invalid path: path escapes workspace")

// ClassifyCommand analyzes the command string and returns its category
func ClassifyCommand(command string) CommandCategory {
	if command == "" {
		return CmdUnknown
	}

	// Normalize command: lowercase for pattern matching
	cmdLower := strings.ToLower(command)

	// First check against dangerous patterns
	for _, pattern := range DangerousPatterns {
		if strings.Contains(cmdLower, strings.ToLower(pattern)) {
			return CmdDangerous
		}
	}

	// Extract base command (first word)
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return CmdUnknown
	}

	baseCmd := parts[0]

	// Check against read-only whitelist
	for _, readonlyCmd := range ReadOnlyCommands {
		if baseCmd == readonlyCmd {
			return CmdReadOnly
		}
	}

	// Check against write commands
	for _, writeCmd := range WriteCommands {
		if baseCmd == writeCmd {
			return CmdWrite
		}
	}

	// Unknown category
	return CmdUnknown
}

// ValidateCommand returns an error if the command is dangerous
func ValidateCommand(command string) error {
	category := ClassifyCommand(command)
	if category == CmdDangerous {
		return errDangerousCommand
	}

	// Check for path injection attempts
	if err := validatePathInjection(command); err != nil {
		return err
	}

	return nil
}

// IsReadOnlyCommand checks if the command is in the read-only whitelist
func IsReadOnlyCommand(command string) bool {
	category := ClassifyCommand(command)
	return category == CmdReadOnly
}

// extractPaths extracts file paths from common commands
func extractPaths(command string) []string {
	var paths []string

	parts := strings.Fields(command)
	if len(parts) < 2 {
		return paths
	}

	// Common patterns for path extraction
	pathIndicators := []string{
		">>", ">", "<", "-i", "-f", "-e", "--",
	}

	for i, part := range parts {
		// Skip flags
		if strings.HasPrefix(part, "-") {
			continue
		}

		// Check if this part might be a path
		isPath := false
		for _, indicator := range pathIndicators {
			if part == indicator && i+1 < len(parts) {
				isPath = true
				break
			}
		}

		// Also check if part looks like a path (contains / or starts with ~)
		if !isPath && (strings.Contains(part, "/") || strings.HasPrefix(part, "~") || strings.HasPrefix(part, ".")) {
			isPath = true
		}

		if isPath && i+1 < len(parts) {
			path := parts[i+1]
			// Skip if it looks like a flag value
			if !strings.HasPrefix(path, "-") {
				paths = append(paths, path)
			}
		}
	}

	return paths
}

// validatePathInjection checks if the command contains path injection attempts
func validatePathInjection(command string) error {
	paths := extractPaths(command)

	for _, path := range paths {
		// Check for path traversal
		if strings.Contains(path, "..") {
			return errPathInjection
		}

		// Check for other dangerous path patterns
		if strings.HasPrefix(path, "/dev/") {
			return errInvalidPath
		}

		// Check for absolute paths that might escape workspace
		// This is a simple check - in production, you'd compare against actual workspace
		if strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "/home/") && !strings.HasPrefix(path, "/tmp/") {
			// Could be system path - be cautious
			continue
		}
	}

	return nil
}

// ValidatePath checks if a given path is safe (no path traversal)
func ValidatePath(path string) error {
	if strings.Contains(path, "..") {
		return errPathInjection
	}

	if strings.HasPrefix(path, "/dev/") {
		return errInvalidPath
	}

	return nil
}

// GetCommandCategoryName returns the name of the category
func GetCommandCategoryName(category CommandCategory) string {
	return category.String()
}
