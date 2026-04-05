package permission

import (
	"path/filepath"
	"regexp"
	"strings"
)

// SemanticValidator provides comprehensive semantic analysis of bash commands
type SemanticValidator struct {
	workingDir string
}

// NewSemanticValidator creates a new SemanticValidator with the specified working directory
func NewSemanticValidator(workingDir string) *SemanticValidator {
	absWorkingDir, _ := filepath.Abs(workingDir)
	return &SemanticValidator{
		workingDir: absWorkingDir,
	}
}

// SemanticReadOnlyCommands is a comprehensive set of read-only commands that don't modify the filesystem
var SemanticReadOnlyCommands = map[string]bool{
	"ls":        true,
	"cat":       true,
	"grep":      true,
	"find":      true,
	"wc":        true,
	"head":      true,
	"tail":      true,
	"echo":      true,
	"pwd":       true,
	"tree":      true,
	"stat":      true,
	"file":      true,
	"du":        true,
	"diff":      true,
	"sort":      true,
	"uniq":      true,
	"cut":       true,
	"tr":        true,
	"man":       true,
	"which":     true,
	"whereis":   true,
	"readlink":  true,
	"basename":  true,
	"dirname":   true,
	"realpath":  true,
	"md5sum":    true,
	"sha256sum": true,
	"sha1sum":   true,
	"cksum":     true,
	"catv":      true,
	"nl":        true,
	"od":        true,
	"xxd":       true,
	"hexdump":   true,
	"printf":    true,
	"date":      true,
	"cal":       true,
	"df":        true,
	"free":      true,
	"uptime":    true,
	"whoami":    true,
	"id":        true,
	"groups":    true,
	"hostname":  true,
	"env":       true,
	"printenv":  true,
	"set":       true,
}

// WriteIndicators are patterns that indicate write operations
var WriteIndicators = []string{
	">",   // Output redirect
	">>",  // Append redirect
	"|",   // Pipe (can lead to write)
	";",   // Command chain
	"&&",  // AND chain
	"||",  // OR chain
	"$()", // Command substitution
	"`",   // Backtick command substitution
	"tee", // tee command
}

// DestructivePatterns is a comprehensive list of dangerous command patterns
var DestructivePatterns = []string{
	"rm -rf",
	"rm -r",
	"rm *",
	"rm -f",
	"rm -i",
	"mv *",
	"mv .",
	"cp -rf",
	"cp -r",
	"cp -f",
	"dd if=",
	"mkfs",
	"fdisk",
	"chmod 777 /",
	"chmod 777 /etc",
	"chmod -R 777",
	"chown root",
	"chown -R",
	"sudo",
	"curl | bash",
	"wget | bash",
	"sh -c",
	"bash -c",
	"> /dev/sda",
	"> /dev/hda",
	"> /dev/nvme",
	":(){:|:&};:",
	"fork bomb",
	"eval ",
	"exec ",
	"chroot",
	"iptables",
	"ufw",
	"firewall-cmd",
	"systemctl stop",
	"systemctl disable",
	"service stop",
	"kill -9 -1",
	"killall",
	"pkill -9",
	"reboot",
	"shutdown",
	"init 0",
	"init 6",
	"halt",
	"poweroff",
}

// BlockedPaths are system paths that should never be accessed
var BlockedPaths = []string{
	"/dev/",
	"/proc/",
	"/sys/",
	"/boot/",
	"/root/",
	"/snap/",
	"/lost+found/",
}

// SedWritePatterns detects sed commands that write to files
var SedWritePatterns = []*regexp.Regexp{
	regexp.MustCompile(`sed\s+(-i|--in-place)\s+['"]?\S+['"]?\s+['"]?([^'"]+)['"]?`),
	regexp.MustCompile(`sed\s+['"]?([^'"]+)['"]?\s*>\s*(\S+)`),
	regexp.MustCompile(`sed\s+['"]?([^'"]+)['"]?\s*>>\s*(\S+)`),
}

// AwkWritePatterns detects awk commands that write to files
var AwkWritePatterns = []*regexp.Regexp{
	regexp.MustCompile(`awk\s+['"]?([^'"]+)\s*>\s*(\S+)`),
	regexp.MustCompile(`awk\s+['"]?([^'"]+)\s*>>\s*(\S+)`),
	regexp.MustCompile(`awk\s+(-f\s+\S+\s+)?['"]?{[^}]*print\s+[^}]*}\s*>\s*(\S+)`),
	regexp.MustCompile(`awk\s+(-f\s+\S+\s+)?['"]?{[^}]*print\s+[^}]*}\s*>>\s*(\S+)`),
}

// ForkBombPattern detects fork bomb attempts
var ForkBombPattern = regexp.MustCompile(`:\s*\(\s*\)\s*\{[^}]*\|[^}]*&[^}]*\};?:`)

// VerifyReadOnly checks if a command is read-only (doesn't modify filesystem)
func (sv *SemanticValidator) VerifyReadOnly(command string) bool {
	if command == "" {
		return false
	}

	// Check for destructive command indicators
	if sv.containsWriteIndicators(command) {
		return false
	}

	// Check for destructive patterns
	if sv.containsDestructivePattern(command) {
		return false
	}

	// Check for sed -i or awk with redirection
	if sv.containsSedWrite(command) || sv.containsAwkWrite(command) {
		return false
	}

	// Extract base command
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return false
	}

	baseCmd := fields[0]

	// Check if base command is in read-only set
	// Also check if the base command itself is a write command
	if !SemanticReadOnlyCommands[baseCmd] {
		return false
	}

	// Additional check: even for read-only commands, ensure no write operations
	// are being performed through the command arguments
	if sv.hasWriteArguments(command) {
		return false
	}

	return true
}

// containsWriteIndicators checks for write operation indicators
func (sv *SemanticValidator) containsWriteIndicators(command string) bool {
	for _, indicator := range WriteIndicators {
		if strings.Contains(command, indicator) {
			return true
		}
	}
	return false
}

// containsDestructivePattern checks for destructive command patterns
func (sv *SemanticValidator) containsDestructivePattern(command string) bool {
	lowerCmd := strings.ToLower(command)
	for _, pattern := range DestructivePatterns {
		if strings.Contains(lowerCmd, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// containsSedWrite checks for sed commands that write to files
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
	// Check for file modification commands as arguments
	writeCommands := []string{"cp", "mv", "rm", "mkdir", "touch", "chmod", "chown", "tee"}

	fields := strings.Fields(command)
	for _, field := range fields {
		for _, wc := range writeCommands {
			if field == wc || strings.HasPrefix(field, wc+" ") || strings.HasPrefix(field, "-"+wc) {
				return true
			}
		}
	}

	return false
}

// DetectDestructive detects if a command is destructive and returns the reason
func (sv *SemanticValidator) DetectDestructive(command string) (bool, string) {
	if command == "" {
		return false, ""
	}

	lowerCmd := strings.ToLower(command)

	// Check for fork bomb
	if ForkBombPattern.MatchString(command) || strings.Contains(lowerCmd, ":(){:|:&}:") {
		return true, "fork bomb detected: recursive function that spawns processes indefinitely"
	}

	// Check each destructive pattern
	destructiveReasons := map[string]string{
		"rm -rf":            "recursive force delete - may delete entire directory tree",
		"rm -r":             "recursive delete - may delete multiple files",
		"rm *":              "wildcard delete - may delete all files in directory",
		"mv *":              "wildcard move - may move all files unexpectedly",
		"cp -rf":            "recursive force copy - may overwrite files",
		"cp -r":             "recursive copy - may copy unintended files",
		"dd if=":            "direct disk access - may overwrite disk data",
		"mkfs":              "filesystem creation - will destroy data",
		"fdisk":             "partition manipulation - may destroy disk data",
		"chmod 777 /":       "world-writable system directory - security vulnerability",
		"chmod 777 /etc":    "world-writable /etc - security vulnerability",
		"chmod -R 777":      "recursive chmod 777 - security vulnerability",
		"chown root":        "ownership change to root - privilege escalation",
		"chown -R":          "recursive ownership change - may change ownership of system files",
		"sudo":              "privilege escalation - executes with elevated permissions",
		"curl | bash":       "remote code execution - executes downloaded script",
		"wget | bash":       "remote code execution - executes downloaded script",
		"sh -c":             "shell execution - may execute arbitrary commands",
		"bash -c":           "shell execution - may execute arbitrary commands",
		"> /dev/sda":        "direct device write - may destroy disk data",
		"> /dev/hda":        "direct device write - may destroy disk data",
		"> /dev/nvme":       "direct device write - may destroy disk data",
		"eval ":             "eval execution - may execute arbitrary commands",
		"exec ":             "exec replacement - may replace current process",
		"chroot":            "chroot operation - may escape to different filesystem",
		"iptables":          "firewall manipulation - may block network access",
		"ufw":               "firewall manipulation - may block network access",
		"firewall-cmd":      "firewall manipulation - may block network access",
		"systemctl stop":    "system service stop - may stop critical services",
		"systemctl disable": "system service disable - may prevent service startup",
		"service stop":      "service stop - may stop critical services",
		"kill -9 -1":        "kill all processes - will terminate all processes",
		"killall":           "kill all processes - will terminate named processes",
		"pkill -9":          "kill processes - will terminate processes",
		"reboot":            "system reboot - will restart the system",
		"shutdown":          "system shutdown - will power off the system",
		"init 0":            "system halt - will halt the system",
		"init 6":            "system reboot - will reboot the system",
		"halt":              "system halt - will halt the system",
		"poweroff":          "system poweroff - will power off the system",
	}

	for pattern, reason := range destructiveReasons {
		if strings.Contains(lowerCmd, strings.ToLower(pattern)) {
			return true, reason
		}
	}

	// Check for dangerous combinations
	if sv.hasDangerousCombination(command) {
		return true, "dangerous command combination detected"
	}

	return false, ""
}

// hasDangerousCombination checks for dangerous command combinations
func (sv *SemanticValidator) hasDangerousCombination(command string) bool {
	// Check for pipe to shell (remote code execution)
	if (strings.Contains(command, "| bash") || strings.Contains(command, "| sh") ||
		strings.Contains(command, "| zsh") || strings.Contains(command, "| fish")) &&
		!strings.HasPrefix(strings.TrimSpace(command), "#") {
		return true
	}

	// Check for Here Document with execute (EOF; bash << EOF)
	if regexp.MustCompile(`<<\s*['"]?EOF['"]?`).MatchString(command) &&
		(strings.Contains(command, "bash") || strings.Contains(command, "sh -c")) {
		return true
	}

	return false
}

// ExtractSedWritePaths extracts file paths from sed write commands
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
				path := match[len(match)-1]
				if path != "" && !strings.HasPrefix(path, "-") {
					paths = append(paths, path)
				}
			}
		}
	}

	// Also check for simple sed redirection patterns
	if strings.Contains(command, "sed ") && (strings.Contains(command, " > ") || strings.Contains(command, " >> ")) {
		re := regexp.MustCompile(`sed\s+[^>]+[>>]?\s+(\S+)`)
		matches := re.FindAllStringSubmatch(command, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				paths = append(paths, match[1])
			}
		}
	}

	return paths
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

	return paths
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

// ExtractAllPaths extracts all file paths from a command
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
	if strings.Contains(path, "..") {
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
	// Use a simple approach - check common environment variables
	switch key {
	case "HOME":
		// Try to get from os.Getenv in real implementation
		// For now, return empty string to trigger fallback
		return ""
	default:
		return ""
	}
}

// ParsePipes parses a command and returns each stage of the pipeline
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

// ParseRedirects parses a command and returns redirect information
func (sv *SemanticValidator) ParseRedirects(command string) []RedirectInfo {
	var redirects []RedirectInfo

	// Patterns for various redirects
	redirectPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(\d?>)\s*(\S+)`),  // > or >file
		regexp.MustCompile(`(\d?>>)\s*(\S+)`), // >> or >>file
		regexp.MustCompile(`(\d*>&?\d+)`),     // 2>&1 style redirects
		regexp.MustCompile(`(\d*<)\s*(\S+)`),  // < input redirect
	}

	// First, isolate the main command from redirects
	// by removing subshells and command substitution
	processed := command
	processed = regexp.MustCompile(`\$\([^)]*\)`).ReplaceAllString(processed, "")
	processed = regexp.MustCompile("`[^`]+`").ReplaceAllString(processed, "")

	for _, pattern := range redirectPatterns {
		matches := pattern.FindAllStringSubmatch(processed, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				redirects = append(redirects, RedirectInfo{
					Type:   match[1],
					Target: match[2],
					FD:     parseFD(match[1]),
				})
			}
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

// ParseSubshells parses subshell and command substitution content
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

// AnalyzeSemantics performs comprehensive semantic analysis on a command
func (sv *SemanticValidator) AnalyzeSemantics(command string) *SemanticAnalysis {
	if command == "" {
		return &SemanticAnalysis{
			IsValid:  false,
			Reason:   "empty command",
			Severity: SeverityError,
		}
	}

	analysis := &SemanticAnalysis{
		Pipes:     sv.ParsePipes(command),
		Redirects: sv.ParseRedirects(command),
		Subshells: sv.ParseSubshells(command),
		Chains:    sv.ParseCommandChaining(command),
	}

	// Check for destructive commands
	isDestructive, reason := sv.DetectDestructive(command)
	if isDestructive {
		analysis.IsValid = false
		analysis.IsDestructive = true
		analysis.Reason = reason
		analysis.Severity = SeverityFatal
		return analysis
	}

	// Validate all paths
	allPaths := sv.ExtractAllPaths(command)
	for _, path := range allPaths {
		valid, pathReason := sv.ValidatePath(path)
		if !valid {
			analysis.IsValid = false
			analysis.InvalidPaths = append(analysis.InvalidPaths, path)
			analysis.Reason = pathReason
			analysis.Severity = maxSeverity(analysis.Severity, SeverityError)
		}
	}

	// Validate sed/awk write paths
	validSedAwk, sedAwkReason := sv.ValidateSedAwkPaths(command)
	if !validSedAwk {
		analysis.IsValid = false
		analysis.Reason = sedAwkReason
		analysis.Severity = maxSeverity(analysis.Severity, SeverityError)
	}

	// Check if command is read-only
	if sv.VerifyReadOnly(command) {
		analysis.IsReadOnly = true
		analysis.IsValid = true
		analysis.Severity = SeverityNone
		return analysis
	}

	// If we have write operations, check if paths are valid
	if len(analysis.Redirects) > 0 {
		for _, redirect := range analysis.Redirects {
			if redirect.FD == 1 || redirect.FD == 3 { // stdout or combined
				valid, reason := sv.ValidatePath(redirect.Target)
				if !valid {
					analysis.IsValid = false
					analysis.InvalidPaths = append(analysis.InvalidPaths, redirect.Target)
					analysis.Reason = reason
					analysis.Severity = maxSeverity(analysis.Severity, SeverityWarning)
				}
			}
		}
	}

	// Check subshells for dangerous content
	for _, subshell := range analysis.Subshells {
		if sv.isDangerousSubshell(subshell.Content) {
			analysis.IsValid = false
			analysis.IsDestructive = true
			analysis.Reason = "dangerous content in subshell"
			analysis.Severity = SeverityFatal
			return analysis
		}
	}

	// If no issues found but command is not read-only, it's a write command (potentially valid)
	if analysis.IsValid && !analysis.IsReadOnly {
		analysis.Severity = SeverityInfo
		analysis.Reason = "write operation"
	}

	return analysis
}

// isDangerousSubshell checks if subshell content is dangerous
func (sv *SemanticValidator) isDangerousSubshell(content string) bool {
	lowerContent := strings.ToLower(content)

	dangerousPatterns := []string{
		"rm -rf",
		"rm -r",
		"curl |",
		"wget |",
		"> /dev/",
		"dd if=",
		"mkfs",
		"fork bomb",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerContent, strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}

// SemanticAnalysis holds the results of semantic analysis
type SemanticAnalysis struct {
	IsValid       bool           // Whether the command is valid
	IsReadOnly    bool           // Whether the command is read-only
	IsDestructive bool           // Whether the command is destructive
	Reason        string         // Reason for invalidity or additional info
	Severity      SeverityLevel  // Severity level of any issues
	InvalidPaths  []string       // List of invalid paths
	Pipes         []string       // Pipeline stages
	Redirects     []RedirectInfo // Redirect information
	Subshells     []SubshellInfo // Subshell information
	Chains        []ChainInfo    // Command chain information
}

// SeverityLevel represents the severity of analysis issues
type SeverityLevel int

const (
	SeverityNone    SeverityLevel = iota // No issues
	SeverityInfo                         // Informational
	SeverityWarning                      // Warning
	SeverityError                        // Error
	SeverityFatal                        // Fatal (command blocked)
)

// maxSeverity returns the higher of two severity levels
func maxSeverity(a, b SeverityLevel) SeverityLevel {
	if a > b {
		return a
	}
	return b
}

// ValidateFullCommand performs complete validation of a command
func (sv *SemanticValidator) ValidateFullCommand(command string) (bool, string, *SemanticAnalysis) {
	analysis := sv.AnalyzeSemantics(command)

	if !analysis.IsValid {
		return false, analysis.Reason, analysis
	}

	if analysis.IsDestructive {
		return false, analysis.Reason, analysis
	}

	if len(analysis.InvalidPaths) > 0 {
		return false, "invalid paths: " + strings.Join(analysis.InvalidPaths, ", "), analysis
	}

	return true, analysis.Reason, analysis
}

// GetWorkingDir returns the working directory of the validator
func (sv *SemanticValidator) GetWorkingDir() string {
	return sv.workingDir
}

// SetWorkingDir sets the working directory of the validator
func (sv *SemanticValidator) SetWorkingDir(workingDir string) {
	sv.workingDir, _ = filepath.Abs(workingDir)
}
