package permission

import (
	"regexp"
	"strings"
)

func (sv *SemanticValidator) ParsePipes(command string) []string {
	// First, handle subshells and command substitution to avoid false positives
	processed := command

	// Remove command substitution content (replace with placeholder)
	processed = regexp.MustCompile(`\$\([^)]*\)`).ReplaceAllString(processed, "$(_)")
	processed = regexp.MustCompile("`[^`]+`").ReplaceAllString(processed, "$(_)")

	// Split by pipe
	stages := strings.Split(processed, "|")

	var result []string
	for _, stage := range stages {
		stage = strings.TrimSpace(stage)
		if stage != "" {
			result = append(result, stage)
		}
	}

	return result
}

func (sv *SemanticValidator) ParseRedirects(command string) []RedirectInfo {
	var redirects []RedirectInfo

	// First, isolate the main command from redirects
	// by removing subshells and command substitution
	processed := command
	processed = regexp.MustCompile(`\$\([^)]*\)`).ReplaceAllString(processed, "")
	processed = regexp.MustCompile("`[^`]+`").ReplaceAllString(processed, "")

	matches := RedirectPattern.FindAllStringSubmatch(processed, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			redirects = append(redirects, RedirectInfo{
				Type:   match[1],
				Target: match[2],
				FD:     parseFD(match[1]),
			})
		}
	}

	return redirects
}

// RedirectInfo holds information about a redirect operation
type RedirectInfo struct {
	Type   string // >, >>, <, 2>, &>, etc.
	Target string // The target file path
	FD     int    // File descriptor (0=stdin, 1=stdout, 2=stderr)
}

// parseFD parses the file descriptor from a redirect type
func parseFD(redirectType string) int {
	switch redirectType {
	case "<":
		return 0 // stdin
	case ">":
		return 1 // stdout
	case ">>":
		return 1 // stdout (append)
	case "2>":
		return 2 // stderr
	case "2>>":
		return 2 // stderr (append)
	case "&>":
		return 3 // combined stdout and stderr
	default:
		return 1 // default to stdout
	}
}

func (sv *SemanticValidator) ParseSubshells(command string) []SubshellInfo {
	var subshells []SubshellInfo

	// Match $() style command substitution
	dollarParen := regexp.MustCompile(`\$\(([^)]+)\)`)
	matches := dollarParen.FindAllStringSubmatchIndex(command, -1)
	for _, match := range matches {
		if len(match) >= 4 {
			content := command[match[2]:match[3]]
			subshells = append(subshells, SubshellInfo{
				Type:    "$()",
				Content: content,
				Start:   match[0],
				End:     match[1],
			})
		}
	}

	// Match backtick command substitution
	backtick := regexp.MustCompile("`([^`]+)`")
	matches = backtick.FindAllStringSubmatchIndex(command, -1)
	for _, match := range matches {
		if len(match) >= 4 {
			content := command[match[2]:match[3]]
			subshells = append(subshells, SubshellInfo{
				Type:    "``",
				Content: content,
				Start:   match[0],
				End:     match[1],
			})
		}
	}

	return subshells
}

// SubshellInfo holds information about a subshell or command substitution
type SubshellInfo struct {
	Type    string // $() or ``
	Content string // The command inside the subshell
	Start   int    // Start position in original command
	End     int    // End position in original command
}

// ParseCommandChaining parses command chaining operators
func (sv *SemanticValidator) ParseCommandChaining(command string) []ChainInfo {
	var chains []ChainInfo

	// First, split by semicolons
	semicolonParts := strings.Split(command, ";")
	currentPos := 0

	for _, part := range semicolonParts {
		part = strings.TrimSpace(part)
		if part == "" {
			currentPos += 1
			continue
		}

		// Check for && and || within this segment
		// Use a regex to split while keeping delimiters
		re := regexp.MustCompile(`(\s*&&\s*|\s*\|\|\s*)`)
		segments := re.Split(part, -1)
		delimiters := re.FindAllString(part, -1)

		for i, segment := range segments {
			segment = strings.TrimSpace(segment)
			if segment == "" {
				continue
			}

			operator := ""
			if i < len(delimiters) {
				operator = strings.TrimSpace(delimiters[i])
			} else if len(segments) > 1 {
				operator = ";" // Default to semicolon for multi-segment
			}

			chains = append(chains, ChainInfo{
				Command:  segment,
				Operator: operator,
				Position: currentPos,
			})

			currentPos += len(segment) + len(operator)
		}

		currentPos += 1 // Account for semicolon
	}

	return chains
}

// ChainInfo holds information about a command in a chain
type ChainInfo struct {
	Command  string // The command part
	Operator string // The operator before this command (; && ||)
	Position int    // Position in the original command
}

func maxSeverity(a, b SeverityLevel) SeverityLevel {
	if a > b {
		return a
	}
	return b
}
