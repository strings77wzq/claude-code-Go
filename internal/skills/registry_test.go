package skills

import (
	"errors"
	"testing"
)

func TestRegistryRegister(t *testing.T) {
	registry := NewRegistry()

	skill := Skill{
		Name:        "test-skill",
		Description: "A test skill",
		Prompt:      "This is a test prompt",
	}

	err := registry.Register(skill)
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}
}

func TestRegistryRegisterDuplicate(t *testing.T) {
	registry := NewRegistry()

	skill1 := Skill{Name: "test", Description: "First", Prompt: "Prompt 1"}
	skill2 := Skill{Name: "test", Description: "Second", Prompt: "Prompt 2"}

	_ = registry.Register(skill1)
	err := registry.Register(skill2)

	if err == nil {
		t.Error("Register duplicate should return error")
	}

	if !errors.Is(err, ErrSkillAlreadyExists) {
		t.Errorf("Expected ErrSkillAlreadyExists, got: %v", err)
	}
}

func TestRegistryRegisterEmptyName(t *testing.T) {
	registry := NewRegistry()

	skill := Skill{
		Name:        "",
		Description: "No name",
		Prompt:      "Prompt",
	}

	err := registry.Register(skill)
	if err == nil {
		t.Error("Register with empty name should return error")
	}

	if !errors.Is(err, ErrInvalidSkill) {
		t.Errorf("Expected ErrInvalidSkill, got: %v", err)
	}
}

func TestRegistryGet(t *testing.T) {
	registry := NewRegistry()

	skill := Skill{
		Name:        "test-skill",
		Description: "Test",
		Prompt:      "Test prompt",
	}

	_ = registry.Register(skill)

	retrieved := registry.Get("test-skill")
	if retrieved == nil {
		t.Fatal("Get should return the skill")
	}

	if retrieved.Name != "test-skill" {
		t.Errorf("Expected name 'test-skill', got '%s'", retrieved.Name)
	}
}

func TestRegistryGetNonExistent(t *testing.T) {
	registry := NewRegistry()

	retrieved := registry.Get("non-existent")
	if retrieved != nil {
		t.Error("Get non-existent skill should return nil")
	}
}

func TestRegistryList(t *testing.T) {
	registry := NewRegistry()

	skill1 := Skill{Name: "skill1", Description: "First", Prompt: "Prompt 1"}
	skill2 := Skill{Name: "skill2", Description: "Second", Prompt: "Prompt 2"}

	_ = registry.Register(skill1)
	_ = registry.Register(skill2)

	list := registry.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 skills, got %d", len(list))
	}
}

func TestRegistryListEmpty(t *testing.T) {
	registry := NewRegistry()

	list := registry.List()
	if len(list) != 0 {
		t.Errorf("Expected empty list, got %d", len(list))
	}
}

func TestRegistryExecute(t *testing.T) {
	registry := NewRegistry()

	skill := Skill{
		Name:        "test-skill",
		Description: "Test",
		Prompt:      "Execute this prompt",
	}

	_ = registry.Register(skill)

	result, err := registry.Execute("test-skill")
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	if result != "Execute this prompt" {
		t.Errorf("Expected prompt, got '%s'", result)
	}
}

func TestRegistryExecuteInvalidName(t *testing.T) {
	registry := NewRegistry()

	result, err := registry.Execute("non-existent")
	if err == nil {
		t.Error("Execute invalid name should return error")
	}

	if !errors.Is(err, ErrSkillNotFound) {
		t.Errorf("Expected ErrSkillNotFound, got: %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty result, got '%s'", result)
	}
}
