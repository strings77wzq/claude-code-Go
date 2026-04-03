package permission

import (
	"fmt"
	"strings"
)

type Rule struct {
	ToolPattern  string
	InputPattern string
}

func ParseRule(ruleString string) (Rule, error) {
	parts := strings.Split(ruleString, "(")
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("invalid rule format: %s", ruleString)
	}

	toolPattern := parts[0]
	inputPart := parts[1]

	if !strings.HasSuffix(inputPart, ")") {
		return Rule{}, fmt.Errorf("invalid rule format: missing closing parenthesis")
	}

	inputPattern := strings.TrimSuffix(inputPart, ")")

	return Rule{
		ToolPattern:  toolPattern,
		InputPattern: inputPattern,
	}, nil
}

func MatchRule(rule Rule, toolName string, input map[string]any) bool {
	if !matchGlob(rule.ToolPattern, toolName) {
		return false
	}

	inputValue := extractInputValue(toolName, input)
	if inputValue == "" {
		return false
	}

	return matchGlob(rule.InputPattern, inputValue)
}

func extractInputValue(toolName string, input map[string]any) string {
	switch toolName {
	case "Bash":
		if cmd, ok := input["command"].(string); ok {
			return cmd
		}
	case "Write", "Edit", "Read":
		if path, ok := input["file_path"].(string); ok {
			return path
		}
	case "Glob", "Grep":
		if pattern, ok := input["pattern"].(string); ok {
			return pattern
		}
	}
	return ""
}

func matchGlob(pattern, text string) bool {
	if pattern == "*" {
		return true
	}

	if pattern == text {
		return true
	}

	patternLower := strings.ToLower(pattern)
	textLower := strings.ToLower(text)

	if patternLower == textLower {
		return true
	}

	parts := strings.SplitN(pattern, ":", 2)
	if len(parts) == 2 {
		prefix := parts[0]
		suffix := parts[1]

		if suffix == "*" {
			if strings.HasPrefix(text, prefix) {
				return true
			}
			if strings.Contains(prefix, "*") {
				return strings.HasSuffix(text, strings.TrimPrefix(prefix, "*"))
			}
			return false
		}

		if prefix == "*" {
			wildcardSuffix := strings.TrimPrefix(suffix, "*")
			if wildcardSuffix == "" {
				return true
			}
			if strings.HasSuffix(text, suffix) {
				return true
			}
			if strings.Contains(suffix, "*") {
				ext := strings.TrimPrefix(suffix, "*")
				return strings.HasSuffix(text, ext)
			}
			return false
		}

		prefixMatch := strings.HasPrefix(text, prefix)
		if !prefixMatch && strings.Contains(prefix, "*") {
			prefixMatch = strings.HasSuffix(text, strings.TrimPrefix(prefix, "*"))
		}

		suffixMatch := strings.HasSuffix(text, suffix)
		if !suffixMatch && strings.Contains(suffix, "*") {
			suffixMatch = strings.HasPrefix(text, strings.TrimSuffix(suffix, "*"))
		}

		return prefixMatch && suffixMatch
	}

	if strings.Contains(pattern, "*") {
		parts2 := strings.SplitN(pattern, "*", 2)
		prefix := parts2[0]

		if prefix != "" && strings.HasPrefix(text, prefix) {
			return true
		}

		if len(parts2) > 1 && parts2[1] != "" {
			suffix := parts2[1]
			return strings.HasSuffix(text, suffix)
		}
	}

	return false
}
