package main

import (
	"fmt"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
)

// PermissionErrorExample demonstrates how to handle permission errors
func main() {
	// Setup permission enforcer
	enforcer := permission.NewEnforcer(
		permission.ModeWorkspaceWrite,
		[]permission.GlobRule{
			{Pattern: "*.go", Allowed: true},
			{Pattern: "*.md", Allowed: true},
			{Pattern: "*.env", Allowed: false}, // Sensitive files
		},
	)

	// Try to read a file
	target := "config.env"
	if err := enforcer.CheckRead(target); err != nil {
		if permErr, ok := err.(*permission.Error); ok {
			fmt.Printf("Permission denied: %s\n", permErr.Message)
			fmt.Println("\nTo allow this action:")
			fmt.Println("1. Add an exception:")
			fmt.Printf("   /allow read %s\n", target)
			fmt.Println("\n2. Switch to a more permissive mode:")
			fmt.Println("   /mode DangerFullAccess")
			fmt.Println("   ⚠️  WARNING: This allows all actions!")
			fmt.Println("\n3. Update your glob rules in ~/.go-code/settings.json")
		}
		return
	}

	fmt.Println("Permission granted!")
}
