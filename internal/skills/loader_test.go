package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSkillsFromValidDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	validSkill := `{
		"name": "test-skill",
		"description": "A test skill",
		"prompt": "This is a test prompt"
	}`

	err := os.WriteFile(filepath.Join(tmpDir, "skill1.json"), []byte(validSkill), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	skills, err := LoadSkills(tmpDir)
	if err != nil {
		t.Errorf("LoadSkills failed: %v", err)
	}

	if len(skills) != 1 {
		t.Errorf("Expected 1 skill, got %d", len(skills))
	}

	if skills[0].Name != "test-skill" {
		t.Errorf("Expected skill name 'test-skill', got '%s'", skills[0].Name)
	}
}

func TestLoadSkillsMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	skill1 := `{"name": "skill1", "description": "First", "prompt": "Prompt 1"}`
	skill2 := `{"name": "skill2", "description": "Second", "prompt": "Prompt 2"}`

	os.WriteFile(filepath.Join(tmpDir, "a.json"), []byte(skill1), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.json"), []byte(skill2), 0644)

	skills, err := LoadSkills(tmpDir)
	if err != nil {
		t.Errorf("LoadSkills failed: %v", err)
	}

	if len(skills) != 2 {
		t.Errorf("Expected 2 skills, got %d", len(skills))
	}
}

func TestLoadSkillsWithWarningsLoadsSkillMarkdown(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "test-driven-development")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	body := `---
name: test-driven-development
description: Use TDD
---

# Test Driven Development

Write a failing test first.
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(body), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := LoadSkillsWithWarnings(tmpDir)
	if err != nil {
		t.Fatalf("LoadSkillsWithWarnings failed: %v", err)
	}
	if len(result.Skills) != 1 {
		t.Fatalf("loaded %d skills, want 1; warnings=%v", len(result.Skills), result.Warnings)
	}
	if result.Skills[0].Name != "test-driven-development" {
		t.Fatalf("name = %q", result.Skills[0].Name)
	}
	if !strings.Contains(result.Skills[0].Prompt, "failing test") {
		t.Fatalf("prompt did not include markdown body: %q", result.Skills[0].Prompt)
	}
}

func TestLoadSkillsInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()

	invalidJSON := `{invalid json content}`

	os.WriteFile(filepath.Join(tmpDir, "bad.json"), []byte(invalidJSON), 0644)

	skills, err := LoadSkills(tmpDir)
	if err != nil {
		t.Errorf("LoadSkills should not fail on bad JSON: %v", err)
	}

	if len(skills) != 0 {
		t.Errorf("Expected 0 skills from invalid JSON, got %d", len(skills))
	}
}

func TestLoadSkillsMalformedJSON(t *testing.T) {
	tmpDir := t.TempDir()

	malformedJSON := `{"name": "test", "prompt": "test"`

	os.WriteFile(filepath.Join(tmpDir, "broken.json"), []byte(malformedJSON), 0644)

	skills, err := LoadSkills(tmpDir)
	if err != nil {
		t.Errorf("LoadSkills should not fail on malformed JSON: %v", err)
	}

	if len(skills) != 0 {
		t.Errorf("Expected 0 skills from malformed JSON, got %d", len(skills))
	}
}

func TestLoadSkillsNonExistentDirectory(t *testing.T) {
	skills, err := LoadSkills("/nonexistent/path/to/skills")
	if err == nil {
		t.Error("LoadSkills should fail on non-existent directory")
	}

	if len(skills) != 0 {
		t.Errorf("Expected 0 skills, got %d", len(skills))
	}
}

func TestLoadSkillsIgnoresNonJSONFiles(t *testing.T) {
	tmpDir := t.TempDir()

	validSkill := `{"name": "test-skill", "description": "Test", "prompt": "Test prompt"}`

	os.WriteFile(filepath.Join(tmpDir, "skill.json"), []byte(validSkill), 0644)
	os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("not a skill"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "data.yaml"), []byte("name: test"), 0644)

	skills, err := LoadSkills(tmpDir)
	if err != nil {
		t.Errorf("LoadSkills failed: %v", err)
	}

	if len(skills) != 1 {
		t.Errorf("Expected 1 skill (JSON only), got %d", len(skills))
	}
}

func TestLoadSkillsIgnoresDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	subdir := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(subdir, 0755)

	validSkill := `{"name": "test-skill", "description": "Test", "prompt": "Test prompt"}`
	os.WriteFile(filepath.Join(tmpDir, "skill.json"), []byte(validSkill), 0644)
	os.WriteFile(filepath.Join(subdir, "nested.json"), []byte(validSkill), 0644)

	skills, err := LoadSkills(tmpDir)
	if err != nil {
		t.Errorf("LoadSkills failed: %v", err)
	}

	if len(skills) != 1 {
		t.Errorf("Expected 1 skill (root level only), got %d", len(skills))
	}
}

func TestLoadSkillsSkipsEmptyName(t *testing.T) {
	tmpDir := t.TempDir()

	noName := `{"description": "No name", "prompt": "Test"}`
	os.WriteFile(filepath.Join(tmpDir, "noname.json"), []byte(noName), 0644)

	skills, err := LoadSkills(tmpDir)
	if err != nil {
		t.Errorf("LoadSkills failed: %v", err)
	}

	if len(skills) != 0 {
		t.Errorf("Expected 0 skills (empty name), got %d", len(skills))
	}
}

func TestLoadSkillsWithWarnings(t *testing.T) {
	tmpDir := t.TempDir()

	validSkill := `{"name": "good", "description": "Valid", "prompt": "OK"}`
	invalidJSON := `{broken`
	noName := `{"description": "No name", "prompt": "Test"}`
	nonJSON := `name: yaml-skill`

	os.WriteFile(filepath.Join(tmpDir, "good.json"), []byte(validSkill), 0644)
	os.WriteFile(filepath.Join(tmpDir, "broken.json"), []byte(invalidJSON), 0644)
	os.WriteFile(filepath.Join(tmpDir, "noname.json"), []byte(noName), 0644)
	os.WriteFile(filepath.Join(tmpDir, "readme.yaml"), []byte(nonJSON), 0644)

	result, err := LoadSkillsWithWarnings(tmpDir)
	if err != nil {
		t.Fatalf("LoadSkillsWithWarnings failed: %v", err)
	}

	if len(result.Skills) != 1 {
		t.Errorf("expected 1 valid skill, got %d", len(result.Skills))
	}

	if result.Skills[0].Name != "good" {
		t.Errorf("expected skill name 'good', got %q", result.Skills[0].Name)
	}

	if len(result.Warnings) != 3 {
		t.Errorf("expected 3 warnings (invalid JSON, empty name, unsupported ext), got %d", len(result.Warnings))
	}

	foundReasons := make(map[string]bool)
	for _, w := range result.Warnings {
		foundReasons[w.Reason] = true
	}
	if !foundReasons["missing required field \"name\""] {
		t.Error("expected warning for missing name field")
	}
}

func TestSkillWarningsDiagnostics(t *testing.T) {
	warnings := []SkillWarning{{
		File:   "bad.json",
		Reason: "invalid JSON: token=secret-token",
	}}

	diagnostics := SkillWarningsDiagnostics(warnings)

	if len(diagnostics) != 1 {
		t.Fatalf("expected one diagnostic, got %#v", diagnostics)
	}
	if diagnostics[0].Component != "skills" || diagnostics[0].Code != "skills.invalid" {
		t.Fatalf("unexpected diagnostic: %#v", diagnostics[0])
	}
	if strings.Contains(diagnostics[0].Format(), "secret-token") {
		t.Fatalf("diagnostic leaked secret: %s", diagnostics[0].Format())
	}
}

func TestLoadSkillsEmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	skills, err := LoadSkills(tmpDir)
	if err != nil {
		t.Errorf("LoadSkills failed on empty directory: %v", err)
	}

	if len(skills) != 0 {
		t.Errorf("Expected 0 skills, got %d", len(skills))
	}
}
