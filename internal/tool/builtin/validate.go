package builtin

import (
	"fmt"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
)

// ResolvePath resolves symlinks and ensures the file path stays within the working directory.
func ResolvePath(filePath, workingDir string) (string, error) {
	resolvedPath, err := permission.ResolveAndValidatePath(filePath, workingDir)
	if err != nil {
		return "", fmt.Errorf("path %q is outside working directory %q", filePath, workingDir)
	}
	return resolvedPath, nil
}

// ValidatePath ensures the file path is within the working directory.
func ValidatePath(filePath, workingDir string) error {
	_, err := ResolvePath(filePath, workingDir)
	return err
}
