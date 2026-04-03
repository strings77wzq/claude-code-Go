package builtin

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidatePath ensures the file path is within the working directory.
func ValidatePath(filePath, workingDir string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}
	absWd, err := filepath.Abs(workingDir)
	if err != nil {
		return fmt.Errorf("failed to resolve working directory: %w", err)
	}
	// Ensure the path starts with the working directory
	// Add separator to prevent prefix matching issues (e.g., /foo vs /foobar)
	if !strings.HasPrefix(absPath+string(filepath.Separator), absWd+string(filepath.Separator)) && absPath != absWd {
		return fmt.Errorf("path %q is outside working directory %q", filePath, workingDir)
	}
	return nil
}
