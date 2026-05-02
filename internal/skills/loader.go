package skills

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/diagnostic"
)

var ErrInvalidSkill = errors.New("invalid skill: missing required field")

// SkillWarning describes a non-fatal issue encountered while loading a skill file.
type SkillWarning struct {
	File   string
	Reason string
}

// LoadResult holds the result of loading skills, including warnings for invalid files.
type LoadResult struct {
	Skills   []Skill
	Warnings []SkillWarning
}

// LoadSkills loads skill files from a directory, skipping invalid files silently.
// Use LoadSkillsWithWarnings for non-fatal validation reporting.
func LoadSkills(dir string) ([]Skill, error) {
	result, err := LoadSkillsWithWarnings(dir)
	if err != nil {
		return nil, err
	}
	return result.Skills, nil
}

// LoadSkillsWithWarnings loads skill files and reports invalid files as warnings.
func LoadSkillsWithWarnings(dir string) (*LoadResult, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := &LoadResult{}
	for _, entry := range entries {
		if entry.IsDir() {
			skillPath := filepath.Join(dir, entry.Name(), "SKILL.md")
			if _, err := os.Stat(skillPath); err == nil {
				loadMarkdownSkill(result, skillPath)
			}
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext != ".json" {
			result.Warnings = append(result.Warnings, SkillWarning{
				File:   entry.Name(),
				Reason: fmt.Sprintf("unsupported file extension %q (only .json supported)", ext),
			})
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			result.Warnings = append(result.Warnings, SkillWarning{
				File:   entry.Name(),
				Reason: fmt.Sprintf("failed to read file: %v", err),
			})
			continue
		}

		var skill Skill
		if err := json.Unmarshal(data, &skill); err != nil {
			result.Warnings = append(result.Warnings, SkillWarning{
				File:   entry.Name(),
				Reason: fmt.Sprintf("invalid JSON: %v", err),
			})
			continue
		}

		if skill.Name == "" {
			result.Warnings = append(result.Warnings, SkillWarning{
				File:   entry.Name(),
				Reason: "missing required field \"name\"",
			})
			continue
		}

		result.Skills = append(result.Skills, skill)
	}

	return result, nil
}

func loadMarkdownSkill(result *LoadResult, path string) {
	data, err := os.ReadFile(path)
	name := filepath.Base(filepath.Dir(path))
	if err != nil {
		result.Warnings = append(result.Warnings, SkillWarning{
			File:   path,
			Reason: fmt.Sprintf("failed to read file: %v", err),
		})
		return
	}
	skill := parseMarkdownSkill(name, string(data))
	if skill.Name == "" {
		result.Warnings = append(result.Warnings, SkillWarning{
			File:   path,
			Reason: "missing required field \"name\"",
		})
		return
	}
	result.Skills = append(result.Skills, skill)
}

func parseMarkdownSkill(fallbackName string, body string) Skill {
	skill := Skill{Name: fallbackName, Prompt: strings.TrimSpace(body)}
	trimmed := strings.TrimSpace(body)
	if !strings.HasPrefix(trimmed, "---") {
		return skill
	}
	parts := strings.SplitN(trimmed, "---", 3)
	if len(parts) < 3 {
		return skill
	}
	frontmatter := parts[1]
	content := strings.TrimSpace(parts[2])
	for _, line := range strings.Split(frontmatter, "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		switch key {
		case "name":
			skill.Name = value
		case "description":
			skill.Description = value
		}
	}
	if content != "" {
		skill.Prompt = content
	}
	return skill
}

func SkillWarningsDiagnostics(warnings []SkillWarning) []diagnostic.Diagnostic {
	diagnostics := make([]diagnostic.Diagnostic, 0, len(warnings))
	for _, warning := range warnings {
		diagnostics = append(diagnostics, diagnostic.Diagnostic{
			Component: "skills",
			Severity:  diagnostic.SeverityWarn,
			Code:      "skills.invalid",
			Summary:   "Invalid skill file",
			Detail:    warning.Reason,
			Metadata: map[string]any{
				"file": warning.File,
			},
		})
	}
	return diagnostics
}
