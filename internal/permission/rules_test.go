package permission

import (
	"testing"
)

func TestParseRuleBashGit(t *testing.T) {
	rule, err := ParseRule("bash(git:*)")
	if err != nil {
		t.Errorf("Failed to parse rule: %v", err)
	}

	if rule.ToolPattern != "bash" {
		t.Errorf("Expected tool pattern 'bash', got '%s'", rule.ToolPattern)
	}

	if rule.InputPattern != "git:*" {
		t.Errorf("Expected input pattern 'git:*', got '%s'", rule.InputPattern)
	}
}

func TestParseRuleFilePath(t *testing.T) {
	rule, err := ParseRule("read(/tmp:*)")
	if err != nil {
		t.Errorf("Failed to parse rule: %v", err)
	}

	if rule.ToolPattern != "read" {
		t.Errorf("Expected tool pattern 'read', got '%s'", rule.ToolPattern)
	}

	if rule.InputPattern != "/tmp:*" {
		t.Errorf("Expected input pattern '/tmp:*', got '%s'", rule.InputPattern)
	}
}

func TestParseRuleInvalidFormat(t *testing.T) {
	_, err := ParseRule("invalid")
	if err == nil {
		t.Error("Expected error for invalid format")
	}
}

func TestParseRuleMissingParenthesis(t *testing.T) {
	_, err := ParseRule("bash(git")
	if err == nil {
		t.Error("Expected error for missing closing parenthesis")
	}
}

func TestMatchRuleBashGitMatchesGitCommit(t *testing.T) {
	rule, _ := ParseRule("bash(git:*)")
	input := map[string]any{"command": "git commit -m \"fix\""}

	matched := MatchRule(rule, "Bash", input)
	if !matched {
		t.Error("Expected rule to match 'git commit' command")
	}
}

func TestMatchRuleBashGitDoesNotMatchRm(t *testing.T) {
	rule, _ := ParseRule("bash(git:*)")
	input := map[string]any{"command": "rm -rf /"}

	matched := MatchRule(rule, "Bash", input)
	if matched {
		t.Error("Expected rule NOT to match 'rm -rf' command")
	}
}

func TestMatchRuleBashRmRFMatchesDangerousRm(t *testing.T) {
	rule, _ := ParseRule("bash(rm -rf:*)")
	input := map[string]any{"command": "rm -rf /"}

	matched := MatchRule(rule, "Bash", input)
	if !matched {
		t.Error("Expected rule to match 'rm -rf' command")
	}
}

func TestMatchRuleBashRmRFDoesNotMatchLs(t *testing.T) {
	rule, _ := ParseRule("bash(rm -rf:*)")
	input := map[string]any{"command": "ls -la"}

	matched := MatchRule(rule, "Bash", input)
	if matched {
		t.Error("Expected rule NOT to match 'ls' command")
	}
}

func TestMatchRuleFilePathRule(t *testing.T) {
	rule, _ := ParseRule("read(/tmp:*)")
	input := map[string]any{"file_path": "/tmp/file.txt"}

	matched := MatchRule(rule, "Read", input)
	if !matched {
		t.Error("Expected rule to match /tmp file path")
	}
}

func TestMatchRuleFilePathRuleNoMatch(t *testing.T) {
	rule, _ := ParseRule("read(/tmp:*)")
	input := map[string]any{"file_path": "/home/file.txt"}

	matched := MatchRule(rule, "Read", input)
	if matched {
		t.Error("Expected rule NOT to match /home file path")
	}
}

func TestMatchRuleWriteToFile(t *testing.T) {
	rule, _ := ParseRule("write(*.txt:*)")
	input := map[string]any{"file_path": "notes.txt"}

	matched := MatchRule(rule, "Write", input)
	if !matched {
		t.Error("Expected rule to match write to .txt file")
	}
}

func TestMatchRuleGlobPattern(t *testing.T) {
	rule, _ := ParseRule("glob(*.go:*)")
	input := map[string]any{"pattern": "*.go"}

	matched := MatchRule(rule, "Glob", input)
	if !matched {
		t.Error("Expected rule to match *.go pattern")
	}
}

func TestMatchRuleGrepPattern(t *testing.T) {
	rule, _ := ParseRule("grep(func:*)")
	input := map[string]any{"pattern": "func main"}

	matched := MatchRule(rule, "Grep", input)
	if !matched {
		t.Error("Expected rule to match grep func pattern")
	}
}

func TestMatchRuleWrongTool(t *testing.T) {
	rule, _ := ParseRule("bash(git:*)")
	input := map[string]any{"command": "git status"}

	matched := MatchRule(rule, "Read", input)
	if matched {
		t.Error("Expected rule NOT to match when tool is wrong")
	}
}

func TestMatchRuleEmptyInput(t *testing.T) {
	rule, _ := ParseRule("bash(git:*)")
	input := map[string]any{}

	matched := MatchRule(rule, "Bash", input)
	if matched {
		t.Error("Expected rule NOT to match when input is empty")
	}
}

func TestMatchRuleWildcardOnly(t *testing.T) {
	rule, _ := ParseRule("bash(*)")
	input := map[string]any{"command": "any command"}

	matched := MatchRule(rule, "Bash", input)
	if !matched {
		t.Error("Expected wildcard rule to match any command")
	}
}

func TestMatchGlobPrefix(t *testing.T) {
	if !matchGlob("git:", "git status") {
		t.Error("Expected git: to match git status")
	}
	if matchGlob("git:", "hg status") {
		t.Error("Expected git: NOT to match hg status")
	}
}

func TestMatchGlobSuffix(t *testing.T) {
	if !matchGlob("*:go", "file.go") {
		t.Error("Expected *:go to match file.go")
	}
	if matchGlob("*:go", "file.txt") {
		t.Error("Expected *:go NOT to match file.txt")
	}
}
