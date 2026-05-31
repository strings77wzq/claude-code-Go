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

// WriteIndicators are patterns that indicate write operations.
// Note: pipe (|), command chain (;), and logical operators (&&, ||) are NOT
// write indicators — they compose commands but do not themselves write.
// Destructive commands in any pipeline stage are still caught by DetectDestructive.
var WriteIndicators = []string{
	">",   // Output redirect
	">>",  // Append redirect
	"$(",  // Command substitution
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

// RemoteScriptPattern detects downloads piped directly into a shell.
var RemoteScriptPattern = regexp.MustCompile(`(?i)\b(curl|wget)\b.*\|\s*(bash|sh|zsh|fish)\b`)

// RedirectPattern detects common file redirection operators and targets.
var RedirectPattern = regexp.MustCompile(`(&>|\d*>>|\d*>|\d*<)\s*(\S+)`)

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

func (sv *SemanticValidator) AnalyzeSemantics(command string) *SemanticAnalysis {
	if command == "" {
		return &SemanticAnalysis{
			IsValid:  false,
			Reason:   "empty command",
			Severity: SeverityError,
		}
	}

	analysis := &SemanticAnalysis{
		IsValid:   true,
		Pipes:     sv.ParsePipes(command),
		Redirects: sv.ParseRedirects(command),
		Subshells: sv.ParseSubshells(command),
		Chains:    sv.ParseCommandChaining(command),
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
		if analysis.IsValid && len(analysis.InvalidPaths) == 0 {
			analysis.Severity = SeverityNone
		}
		return analysis
	}

	// If we have write operations, check if paths are valid
	if len(analysis.Redirects) > 0 {
		for _, redirect := range analysis.Redirects {
			if redirect.FD == 1 || redirect.FD == 2 || redirect.FD == 3 { // stdout, stderr, or combined
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

	// Check for destructive commands
	isDestructive, reason := sv.DetectDestructive(command)
	if isDestructive {
		analysis.IsValid = false
		analysis.IsDestructive = true
		analysis.Reason = reason
		analysis.Severity = SeverityFatal
		return analysis
	}

	// If no issues found but command is not read-only, it's a write command (potentially valid)
	if analysis.IsValid && !analysis.IsReadOnly {
		analysis.Severity = SeverityInfo
		analysis.Reason = "write operation"
	}

	return analysis
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
