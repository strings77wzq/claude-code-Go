package permission

import (
	"regexp"
	"strings"
)

func (sv *SemanticValidator) containsDestructivePattern(command string) bool {
	lowerCmd := strings.ToLower(command)
	for _, pattern := range DestructivePatterns {
		if strings.Contains(lowerCmd, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func (sv *SemanticValidator) DetectDestructive(command string) (bool, string) {
	if command == "" {
		return false, ""
	}

	lowerCmd := strings.ToLower(command)

	// Check for fork bomb
	if ForkBombPattern.MatchString(command) || strings.Contains(lowerCmd, ":(){:|:&}:") {
		return true, "fork bomb detected: recursive function that spawns processes indefinitely"
	}

	if RemoteScriptPattern.MatchString(command) {
		return true, "remote code execution - executes downloaded script"
	}

	// Check each destructive pattern in priority order so more specific
	// patterns such as "rm -rf" are not shadowed by broader ones.
	destructiveReasons := []struct {
		pattern string
		reason  string
	}{
		{"rm -rf", "recursive force delete - may delete entire directory tree"},
		{"rm -r", "recursive delete - may delete multiple files"},
		{"rm *", "wildcard delete - may delete all files in directory"},
		{"mv *", "wildcard move - may move all files unexpectedly"},
		{"cp -rf", "recursive force copy - may overwrite files"},
		{"cp -r", "recursive copy - may copy unintended files"},
		{"dd if=", "direct disk access - may overwrite disk data"},
		{"mkfs", "filesystem creation - will destroy data"},
		{"fdisk", "partition manipulation - may destroy disk data"},
		{"chmod 777 /", "world-writable system directory - security vulnerability"},
		{"chmod 777 /etc", "world-writable /etc - security vulnerability"},
		{"chmod -R 777", "recursive chmod 777 - security vulnerability"},
		{"chown root", "ownership change to root - privilege escalation"},
		{"chown -R", "recursive ownership change - may change ownership of system files"},
		{"sudo", "privilege escalation - executes with elevated permissions"},
		{"curl | bash", "remote code execution - executes downloaded script"},
		{"wget | bash", "remote code execution - executes downloaded script"},
		{"sh -c", "shell execution - may execute arbitrary commands"},
		{"bash -c", "shell execution - may execute arbitrary commands"},
		{"> /dev/sda", "direct device write - may destroy disk data"},
		{"> /dev/hda", "direct device write - may destroy disk data"},
		{"> /dev/nvme", "direct device write - may destroy disk data"},
		{"eval ", "eval execution - may execute arbitrary commands"},
		{"exec ", "exec replacement - may replace current process"},
		{"chroot", "chroot operation - may escape to different filesystem"},
		{"iptables", "firewall manipulation - may block network access"},
		{"ufw", "firewall manipulation - may block network access"},
		{"firewall-cmd", "firewall manipulation - may block network access"},
		{"systemctl stop", "system service stop - may stop critical services"},
		{"systemctl disable", "system service disable - may prevent service startup"},
		{"service stop", "service stop - may stop critical services"},
		{"kill -9 -1", "kill all processes - will terminate all processes"},
		{"killall", "kill all processes - will terminate named processes"},
		{"pkill -9", "kill processes - will terminate processes"},
		{"reboot", "system reboot - will restart the system"},
		{"shutdown", "system shutdown - will power off the system"},
		{"init 0", "system halt - will halt the system"},
		{"init 6", "system reboot - will reboot the system"},
		{"halt", "system halt - will halt the system"},
		{"poweroff", "system poweroff - will power off the system"},
	}

	for _, item := range destructiveReasons {
		if strings.Contains(lowerCmd, strings.ToLower(item.pattern)) {
			return true, item.reason
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
