package skills

import (
	"errors"
	"sync"
)

var ErrSkillNotFound = errors.New("skill not found")
var ErrSkillAlreadyExists = errors.New("skill already exists")

type Registry struct {
	mu     sync.RWMutex
	skills map[string]Skill
}

func NewRegistry() *Registry {
	return &Registry{
		skills: make(map[string]Skill),
	}
}

func (r *Registry) Register(skill Skill) error {
	if skill.Name == "" {
		return ErrInvalidSkill
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.skills[skill.Name]; exists {
		return ErrSkillAlreadyExists
	}
	r.skills[skill.Name] = skill
	return nil
}

func (r *Registry) Get(name string) *Skill {
	r.mu.RLock()
	defer r.mu.RUnlock()
	skill, ok := r.skills[name]
	if !ok {
		return nil
	}
	return &skill
}

func (r *Registry) List() []Skill {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Skill, 0, len(r.skills))
	for _, skill := range r.skills {
		result = append(result, skill)
	}
	return result
}

func (r *Registry) Execute(name string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	skill, ok := r.skills[name]
	if !ok {
		return "", ErrSkillNotFound
	}
	return skill.Prompt, nil
}
