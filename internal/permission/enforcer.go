package permission

import "strings"

// PermissionLevel represents the required permission level for a tool
type PermissionLevel int

const (
	// LevelReadOnly allows only read operations
	LevelReadOnly PermissionLevel = iota
	// LevelWorkspaceWrite allows read and write within workspace
	LevelWorkspaceWrite
	// LevelDangerFullAccess allows all operations without restrictions
	LevelDangerFullAccess
)

// String returns the string representation of PermissionLevel
func (p PermissionLevel) String() string {
	switch p {
	case LevelReadOnly:
		return "ReadOnly"
	case LevelWorkspaceWrite:
		return "WorkspaceWrite"
	case LevelDangerFullAccess:
		return "DangerFullAccess"
	default:
		return "Unknown"
	}
}

// Enforcer enforces permission policies with additional tool-level checks
type Enforcer struct {
	policy              *Policy
	bashValidator       func(command string) error
	fileBoundaryChecker func(filePath, workingDir string) (string, error)
}

// NewEnforcer creates a new Enforcer with the given policy
func NewEnforcer(policy *Policy) *Enforcer {
	return &Enforcer{
		policy:              policy,
		bashValidator:       ValidateCommand,
		fileBoundaryChecker: ResolveAndValidatePath,
	}
}

// Evaluate determines the permission decision for a tool execution
// Evaluation chain:
// 1. If tool is Bash: run bash validation → if dangerous, return Deny; if readonly, return Allow
// 2. If tool is Write/Edit: run file boundary check → if escapes workspace, return Deny
// 3. Call original policy.Evaluate()
func (e *Enforcer) Evaluate(toolName string, toolRequiresPermission bool, input map[string]any, workingDir string) Decision {
	// Step 1: Bash tool validation
	if toolName == "Bash" {
		if cmd, ok := input["command"].(string); ok && cmd != "" {
			// Run bash validation
			if err := e.bashValidator(cmd); err != nil {
				return Deny
			}
			// If readonly command, allow without further checks
			if IsReadOnlyCommand(cmd) {
				return Allow
			}
		}
	}

	// Step 2: Write/Edit tool file boundary check
	if toolName == "Write" || toolName == "Edit" {
		if filePath, ok := input["file_path"].(string); ok && filePath != "" && workingDir != "" {
			// Run file boundary check
			if _, err := e.fileBoundaryChecker(filePath, workingDir); err != nil {
				return Deny
			}
		}
	}

	// Step 3: Fall back to original policy evaluation
	return e.policy.Evaluate(toolName, input, toolRequiresPermission)
}

// ValidateBashCommand is a convenience function for validating bash commands
func ValidateBashCommand(command string) error {
	return ValidateCommand(command)
}

// CheckFileBoundary validates that a file path is within the working directory
func CheckFileBoundary(filePath, workingDir string) (string, error) {
	return ResolveAndValidatePath(filePath, workingDir)
}

// IsReadOnlyTool checks if a tool only requires read-only permissions
func IsReadOnlyTool(toolName string) bool {
	readonlyTools := map[string]bool{
		"Read": true,
		"Glob": true,
		"Grep": true,
		"Diff": true,
		"Tree": true,
	}
	return readonlyTools[toolName]
}

// GetToolPermissionLevel returns the default permission level for a tool
func GetToolPermissionLevel(toolName string) PermissionLevel {
	switch toolName {
	case "Read", "Glob", "Grep", "Diff", "Tree", "WebFetch":
		return LevelReadOnly
	case "Write", "Edit", "TodoWrite":
		return LevelWorkspaceWrite
	case "Bash":
		return LevelDangerFullAccess
	default:
		// Check if it's a readonly tool by name pattern
		if IsReadOnlyTool(toolName) {
			return LevelReadOnly
		}
		return LevelWorkspaceWrite
	}
}

// ModeToPermissionLevel converts a Mode to PermissionLevel
func ModeToPermissionLevel(mode Mode) PermissionLevel {
	switch mode {
	case ReadOnly:
		return LevelReadOnly
	case WorkspaceWrite:
		return LevelWorkspaceWrite
	case DangerFullAccess:
		return LevelDangerFullAccess
	default:
		return LevelReadOnly
	}
}

// PermissionLevelToMode converts a PermissionLevel to Mode
func PermissionLevelToMode(level PermissionLevel) Mode {
	switch level {
	case LevelReadOnly:
		return ReadOnly
	case LevelWorkspaceWrite:
		return WorkspaceWrite
	case LevelDangerFullAccess:
		return DangerFullAccess
	default:
		return ReadOnly
	}
}

// ContainsPathTraversal checks if input contains path traversal patterns
func ContainsPathTraversal(input map[string]any) bool {
	// Check common path fields
	pathFields := []string{"file_path", "path", "directory", "folder"}

	for _, field := range pathFields {
		if value, ok := input[field].(string); ok {
			if strings.Contains(value, "..") {
				return true
			}
		}
	}

	return false
}
