package skills

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var ErrInvalidSkill = errors.New("invalid skill: missing required field")

func LoadSkills(dir string) ([]Skill, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var skills []Skill
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var skill Skill
		if err := json.Unmarshal(data, &skill); err != nil {
			continue
		}

		if skill.Name == "" {
			continue
		}

		skills = append(skills, skill)
	}

	return skills, nil
}
