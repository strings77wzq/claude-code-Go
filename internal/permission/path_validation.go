package permission

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (sv *SemanticValidator) ExtractAllPaths(command string) []string {
	var paths []string

	// Skip comments
	if strings.HasPrefix(strings.TrimSpace(command), "#") {
		return paths
	}

	// Extract paths from various patterns
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?:^|\s)([./][^\s]+)|(?:^|\s)(/[^\s]+)`),
		regexp.MustCompile(`[>>]?\s*(\S+[^\s;|&&]*)`),
		regexp.MustCompile(`(?:cat|ls|grep|find|stat|file|du|diff|head|tail|wc|cp|mv|rm|mkdir|touch|chmod|chown)\s+([^\s;|&&]+)`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(command, -1)
		for _, match := range matches {
			for i := 1; i < len(match); i++ {
				if match[i] != "" {
					path := strings.Trim(match[i], "'\"")
					if sv.looksLikePath(path) && !strings.HasPrefix(path, "-") {
						paths = append(paths, path)
					}
				}
			}
		}
	}

	// Deduplicate paths
	seen := make(map[string]bool)
	var uniquePaths []string
	for _, p := range paths {
		if !seen[p] {
			seen[p] = true
			uniquePaths = append(uniquePaths, p)
		}
	}

	return uniquePaths
}

// looksLikePath checks if a string looks like a file path
func (sv *SemanticValidator) looksLikePath(s string) bool {
	if s == "" {
		return false
	}
	// Check for common path indicators
	return strings.HasPrefix(s, "/") ||
		strings.HasPrefix(s, "./") ||
		strings.HasPrefix(s, "../") ||
		strings.HasPrefix(s, "~") ||
		strings.Contains(s, "/") && !strings.Contains(s, "://")
}

// ValidatePath validates that a path is within the working directory and not blocked
func (sv *SemanticValidator) ValidatePath(path string) (bool, string) {
	if path == "" {
		return false, "empty path"
	}

	originalPath := path

	// Expand home directory
	if strings.HasPrefix(path, "~/") {
		home := getHomeDir()
		if home != "" {
			path = filepath.Join(home, path[2:])
		}
	}

	// Handle relative paths
	if !filepath.IsAbs(path) {
		if sv.workingDir != "" {
			path = filepath.Join(sv.workingDir, path)
		}
	}

	// Resolve symlinks
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		// If we can't resolve symlinks, use the original path
		// but still validate the original path
		resolvedPath = path
	}

	// Normalize the path
	resolvedPath = filepath.Clean(resolvedPath)

	// Check for blocked paths
	for _, blocked := range BlockedPaths {
		if strings.HasPrefix(resolvedPath, blocked) {
			return false, "blocked system path: " + blocked
		}
	}

	// Check for path traversal attempts
	if strings.Contains(originalPath, "..") {
		// Check if the resolved path escapes working directory
		if sv.workingDir != "" {
			absWorkingDir := filepath.Clean(sv.workingDir)
			if !strings.HasPrefix(resolvedPath+string(filepath.Separator), absWorkingDir+string(filepath.Separator)) {
				return false, "path traversal attempt detected"
			}
		}
	}

	// If working directory is set, validate path is within it
	if sv.workingDir != "" {
		absWorkingDir := filepath.Clean(sv.workingDir)
		if !strings.HasPrefix(resolvedPath+string(filepath.Separator), absWorkingDir+string(filepath.Separator)) {
			return false, "path escapes working directory: " + absWorkingDir
		}
	}

	return true, ""
}

// getHomeDir returns the user's home directory
func getHomeDir() string {
	// Try to get home directory from environment
	home := getEnv("HOME")
	if home != "" {
		return home
	}

	// Fallback to /tmp for safety
	return "/tmp"
}

// getEnv returns an environment variable value
func getEnv(key string) string {
	return os.Getenv(key)
}
