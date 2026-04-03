package tty

import (
	"fmt"
	"os"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

type Renderer struct{}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) PrintWelcome(version string) {
	fmt.Println(ColorGreen + `
  ____   _    ____ ___ 
 |  _ \ / \  / ___|_ _|
 | |_) / _ \ \___ \| | 
 |  __/ ___ \ ___) | | 
 |_| /_/   \_\____/___|
` + ColorReset)
	fmt.Printf(ColorGreen+"Welcome to go-code %s\n"+ColorReset, version)
	fmt.Println("Type /help for available commands")
	fmt.Println()
}

func (r *Renderer) PrintPrompt() {
	fmt.Print(ColorGreen + "go-code> " + ColorReset)
}

func (r *Renderer) PrintStreaming(text string) {
	fmt.Print(text)
}

func (r *Renderer) PrintToolCall(name string, input map[string]any) {
	fmt.Print(ColorYellow)
	fmt.Printf("[Tool] %s", name)
	fmt.Println(ColorReset)
}

func (r *Renderer) PrintToolResult(name string, result string) {
	fmt.Print(ColorCyan)
	fmt.Printf("[Result] %s: ", name)
	fmt.Println(result)
	fmt.Print(ColorReset)
}

func (r *Renderer) PrintError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, ColorRed+"Error: "+err.Error()+ColorReset)
}

func (r *Renderer) PrintHelp() {
	fmt.Println(ColorCyan + "Available commands:" + ColorReset)
	fmt.Println("  /help   - Show this help message")
	fmt.Println("  /clear  - Clear conversation history")
	fmt.Println("  /exit   - Exit the program")
	fmt.Println("  /quit   - Exit the program")
	fmt.Println("  /model  - Show current model")
}

func (r *Renderer) PrintModel(model string) {
	fmt.Printf(ColorCyan+"Current model: %s\n"+ColorReset, model)
}
