package skills

// Skill represents a custom command skill that can be invoked by the user.
// Skills are essentially named prompts that get injected into the agent's
// system prompt when invoked via /<skill-name> in the REPL.
type Skill struct {
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	Prompt      string   `json:"prompt" yaml:"prompt"`
	Examples    []string `json:"examples" yaml:"examples"`
}
