package permission

// Mode represents the 3-tier permission mode
type Mode string

const (
	// ReadOnly allows only read operations (Read, Glob, Grep)
	ReadOnly Mode = "ReadOnly"
	// WorkspaceWrite allows read and write operations within workspace
	WorkspaceWrite Mode = "WorkspaceWrite"
	// DangerFullAccess allows all operations without restrictions
	DangerFullAccess Mode = "DangerFullAccess"
)

// Decision represents the permission decision
type Decision string

const (
	// Allow permits the operation
	Allow Decision = "Allow"
	// Deny blocks the operation
	Deny Decision = "Deny"
	// Ask prompts the user for a decision
	Ask Decision = "Ask"
)

// Policy defines the permission policy with rules and mode settings
type Policy struct {
	activeMode       Mode
	toolRequirements map[string]Mode
	allowRules       []Rule
	denyRules        []Rule
	sessionMemory    map[string]Decision
}

// NewPolicy creates a new Policy with the specified mode
func NewPolicy(mode Mode) *Policy {
	return &Policy{
		activeMode:       mode,
		toolRequirements: make(map[string]Mode),
		allowRules:       make([]Rule, 0),
		denyRules:        make([]Rule, 0),
		sessionMemory:    make(map[string]Decision),
	}
}

// SetToolRequirement sets the minimum permission mode required for a specific tool
func (p *Policy) SetToolRequirement(tool string, mode Mode) {
	p.toolRequirements[tool] = mode
}

// AddAllowRule adds an allow rule to the policy
func (p *Policy) AddAllowRule(rule Rule) {
	p.allowRules = append(p.allowRules, rule)
}

// AddDenyRule adds a deny rule to the policy
func (p *Policy) AddDenyRule(rule Rule) {
	p.denyRules = append(p.denyRules, rule)
}

// SetSessionMemory stores a decision in session memory
func (p *Policy) SetSessionMemory(key string, decision Decision) {
	p.sessionMemory[key] = decision
}

// GetSessionMemory retrieves a decision from session memory
func (p *Policy) GetSessionMemory(key string) (Decision, bool) {
	decision, exists := p.sessionMemory[key]
	return decision, exists
}

// GetActiveMode returns the current active mode
func (p *Policy) GetActiveMode() Mode {
	return p.activeMode
}

// Evaluate determines the permission decision for a tool execution
// Evaluation priority: deny rules > allow rules > session memory > tool attribute > default Ask
func (p *Policy) Evaluate(toolName string, input map[string]any, requiresPermission bool) Decision {
	memKey := p.generateMemoryKey(toolName, input)

	// 1. Check deny rules (highest priority)
	for _, rule := range p.denyRules {
		if MatchRule(rule, toolName, input) {
			return Deny
		}
	}

	// 2. Check allow rules
	for _, rule := range p.allowRules {
		if MatchRule(rule, toolName, input) {
			return Allow
		}
	}

	// 3. Check session memory (user's "always" choices)
	if decision, exists := p.GetSessionMemory(memKey); exists {
		return decision
	}

	// 4. Check tool requirement - overrides mode, must be explicitly allowed or denied
	if requiredMode, hasReq := p.toolRequirements[toolName]; hasReq {
		if p.activeMode == requiredMode {
			return Allow
		}
		return Deny
	}

	// 5. Check requiresPermission flag
	if !requiresPermission {
		return Allow
	}

	// 6. Mode-based permissions
	if p.activeMode == DangerFullAccess {
		return Allow
	}

	// Default: ask user
	return Ask
}

// meetsModeRequirement checks if the current mode meets the required mode
func (p *Policy) meetsModeRequirement(required Mode) bool {
	modeHierarchy := map[Mode]int{
		ReadOnly:         1,
		WorkspaceWrite:   2,
		DangerFullAccess: 3,
	}

	currentLevel := modeHierarchy[p.activeMode]
	requiredLevel := modeHierarchy[required]

	return currentLevel >= requiredLevel
}

// generateMemoryKey creates a unique key for session memory based on tool and input
func (p *Policy) generateMemoryKey(toolName string, input map[string]any) string {
	// Create a simple key based on tool name and relevant input fields
	var key string

	switch toolName {
	case "Bash":
		if cmd, ok := input["command"].(string); ok {
			key = "bash:" + cmd
		} else {
			key = "bash:unknown"
		}
	case "Read", "Write", "Edit":
		if path, ok := input["file_path"].(string); ok {
			key = toolName + ":" + path
		} else {
			key = toolName + ":unknown"
		}
	case "Glob", "Grep":
		if pattern, ok := input["pattern"].(string); ok {
			key = toolName + ":" + pattern
		} else {
			key = toolName + ":unknown"
		}
	default:
		key = toolName
	}

	return key
}
