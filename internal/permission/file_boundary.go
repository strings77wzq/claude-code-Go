package permission

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// BinaryExtensions maps known binary file extensions
var BinaryExtensions = map[string]bool{
	".exe":   true,
	".bin":   true,
	".so":    true,
	".dylib": true,
	".dll":   true,
	".o":     true,
	".a":     true,
	".lib":   true,
	".pyc":   true,
	".pyo":   true,
	".class": true,
	".jar":   true,
}

// MaxFileSize is the maximum allowed file size (10MB)
const MaxFileSize = 10 * 1024 * 1024

// ErrFileTooLarge is returned when a file exceeds MaxFileSize
var ErrFileTooLarge = errors.New("file exceeds maximum allowed size")

// ErrPathEscape is returned when a symlink resolves outside working directory
var ErrPathEscape = errors.New("path escapes working directory")

// IsBinaryFile checks if a file is binary based on its extension
// For unknown extensions, it reads the first 512 bytes to detect null bytes
func IsBinaryFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	if BinaryExtensions[ext] {
		return true
	}

	// For unknown extensions, check for binary content
	file, err := os.Open(filePath)
	if err != nil {
		// If we can't open the file, assume it's not binary
		// This handles the case where the file doesn't exist yet
		return false
	}
	defer file.Close()

	// Read first 512 bytes to check for null bytes
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil || n == 0 {
		// If we can't read the file or it's empty, assume it's not binary
		return false
	}

	// Check for null bytes (binary indicator)
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return true
		}
	}

	return false
}

// CheckFileSize returns the file size or an error if it exceeds MaxFileSize
func CheckFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		// If the file doesn't exist, that's fine - we'll check before writing
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	if info.Size() > MaxFileSize {
		return info.Size(), ErrFileTooLarge
	}

	return info.Size(), nil
}

// ResolveAndValidatePath resolves symlinks and validates the path is within workingDir
func ResolveAndValidatePath(filePath, workingDir string) (string, error) {
	absWorkingDir, err := filepath.Abs(workingDir)
	if err != nil {
		return "", err
	}
	if resolvedWorkingDir, err := filepath.EvalSymlinks(absWorkingDir); err == nil {
		absWorkingDir = resolvedWorkingDir
	}

	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(absWorkingDir, filePath)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	resolvedPath, err := resolvePathForExistingParent(absPath)
	if err != nil {
		return "", err
	}

	if !isPathWithin(resolvedPath, absWorkingDir) {
		return "", ErrPathEscape
	}

	return resolvedPath, nil
}

func resolvePathForExistingParent(absPath string) (string, error) {
	if resolvedPath, err := filepath.EvalSymlinks(absPath); err == nil {
		return filepath.Clean(resolvedPath), nil
	}

	var missing []string
	current := absPath
	for {
		if resolvedParent, err := filepath.EvalSymlinks(current); err == nil {
			for i := len(missing) - 1; i >= 0; i-- {
				resolvedParent = filepath.Join(resolvedParent, missing[i])
			}
			return filepath.Clean(resolvedParent), nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return filepath.Clean(absPath), nil
		}
		missing = append(missing, filepath.Base(current))
		current = parent
	}
}

func isPathWithin(path, root string) bool {
	path = filepath.Clean(path)
	root = filepath.Clean(root)
	if path == root {
		return true
	}
	return strings.HasPrefix(path+string(filepath.Separator), root+string(filepath.Separator))
}
