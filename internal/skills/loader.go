package skills

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
