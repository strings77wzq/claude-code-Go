package permission

import (
	"strings"
)

func (sv *SemanticValidator) containsWriteIndicators(command string) bool {
	for _, indicator := range WriteIndicators {
		if strings.Contains(command, indicator) {
			return true
		}
	}
	return false
}

func (sv *SemanticValidator) containsSedWrite(command string) bool {
	// Check for sed -i
	if strings.Contains(command, "sed") {
		for _, pattern := range SedWritePatterns {
			if pattern.MatchString(command) {
				return true
			}
		}
	}
	return false
}

// containsAwkWrite checks for awk commands that write to files
func (sv *SemanticValidator) containsAwkWrite(command string) bool {
	if strings.Contains(command, "awk") {
		for _, pattern := range AwkWritePatterns {
			if pattern.MatchString(command) {
				return true
			}
		}
	}
	return false
}

// hasWriteArguments checks for write operations in command arguments
func (sv *SemanticValidator) hasWriteArguments(command string) bool {
	// Check for file modification commands as arguments.
	// Only exact-match: fields are whitespace-split so no field contains a space,
	// and prefix checks like "-cp" are not meaningful write-command indicators.
	writeCommands := []string{"cp", "mv", "rm", "mkdir", "touch", "chmod", "chown", "tee"}

	fields := strings.Fields(command)
	for _, field := range fields {
		for _, wc := range writeCommands {
			if field == wc {
				return true
			}
		}
	}

	return false
}

func (sv *SemanticValidator) ExtractSedWritePaths(command string) []string {
	var paths []string

	if !strings.Contains(command, "sed") {
		return paths
	}

	for _, pattern := range SedWritePatterns {
		matches := pattern.FindAllStringSubmatch(command, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				// The last match group is the file path
				path := strings.TrimSpace(match[len(match)-1])
				parts := strings.FieldsFunc(path, func(r rune) bool {
					return r == ';' || r == '&' || r == '|'
				})
				if len(parts) > 0 {
					path = strings.TrimSpace(parts[0])
				}
				if path != "" && !strings.HasPrefix(path, "-") {
					paths = append(paths, path)
				}
			}
		}
	}

	// Also check for simple sed redirection patterns
	if strings.Contains(command, "sed ") && (strings.Contains(command, " > ") || strings.Contains(command, " >> ")) {
		for _, redirect := range sv.ParseRedirects(command) {
			if redirect.FD == 1 || redirect.FD == 3 {
				paths = append(paths, redirect.Target)
			}
		}
	}

	return uniqueNonFlagPaths(paths)
}

// ExtractAwkWritePaths extracts file paths from awk write commands
func (sv *SemanticValidator) ExtractAwkWritePaths(command string) []string {
	var paths []string

	if !strings.Contains(command, "awk") {
		return paths
	}

	for _, pattern := range AwkWritePatterns {
		matches := pattern.FindAllStringSubmatch(command, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				path := match[len(match)-1]
				if path != "" && !strings.HasPrefix(path, "-") {
					paths = append(paths, path)
				}
			}
		}
	}

	for _, redirect := range sv.ParseRedirects(command) {
		if redirect.FD == 1 || redirect.FD == 3 {
			paths = append(paths, redirect.Target)
		}
	}

	return uniqueNonFlagPaths(paths)
}

func uniqueNonFlagPaths(paths []string) []string {
	seen := make(map[string]bool, len(paths))
	unique := make([]string, 0, len(paths))
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path == "" || strings.HasPrefix(path, "-") || seen[path] {
			continue
		}
		seen[path] = true
		unique = append(unique, path)
	}
	return unique
}

// ValidateSedAwkPaths validates that sed/awk write targets are within workspace
func (sv *SemanticValidator) ValidateSedAwkPaths(command string) (bool, string) {
	sedPaths := sv.ExtractSedWritePaths(command)
	awkPaths := sv.ExtractAwkWritePaths(command)

	allPaths := append(sedPaths, awkPaths...)

	for _, path := range allPaths {
		valid, reason := sv.ValidatePath(path)
		if !valid {
			return false, reason
		}
	}

	return true, ""
}
