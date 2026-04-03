package permission

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Prompter interface {
	Decide(toolName string, input map[string]any, reason string) Decision
}

type Reader interface {
	ReadString(delim byte) (string, error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type TerminalPrompter struct {
	reader Reader
	writer Writer
}

func NewTerminalPrompter(reader Reader, writer Writer) *TerminalPrompter {
	return &TerminalPrompter{
		reader: reader,
		writer: writer,
	}
}

func (p *TerminalPrompter) Decide(toolName string, input map[string]any, reason string) Decision {
	detail := p.formatDetail(toolName, input)

	p.printPrompt(toolName, detail)

	for {
		p.writer.Write([]byte("Allow? (y)es / (n)o / (a)lways: "))

		line, err := p.reader.ReadString('\n')
		if err != nil {
			return Deny
		}

		line = strings.TrimSpace(strings.ToLower(line))

		switch line {
		case "y", "yes":
			return Allow
		case "n", "no":
			return Deny
		case "a", "always":
			return Allow
		}
	}
}

func (p *TerminalPrompter) formatDetail(toolName string, input map[string]any) string {
	switch toolName {
	case "Bash":
		if cmd, ok := input["command"].(string); ok {
			return fmt.Sprintf("Command: %s", cmd)
		}
	case "Write", "Edit":
		if path, ok := input["file_path"].(string); ok {
			return fmt.Sprintf("File: %s", path)
		}
	case "Read":
		if path, ok := input["file_path"].(string); ok {
			return fmt.Sprintf("File: %s", path)
		}
	case "Glob":
		if pattern, ok := input["pattern"].(string); ok {
			return fmt.Sprintf("Pattern: %s", pattern)
		}
	case "Grep":
		if pattern, ok := input["pattern"].(string); ok {
			return fmt.Sprintf("Pattern: %s", pattern)
		}
	}
	return "No details available"
}

func (p *TerminalPrompter) printPrompt(toolName string, detail string) {
	border := "┌──────────────────────────────────────────────┐"
	inner := "│  Tool: %-37s │"
	detailLine := "│  %-44s │"
	spacer := "│                                              │"
	promptLine := "│  %-44s │"

	p.writer.Write([]byte(border + "\n"))
	p.writer.Write([]byte(fmt.Sprintf(inner, toolName) + "\n"))
	p.writer.Write([]byte(spacer + "\n"))
	p.writer.Write([]byte(fmt.Sprintf(detailLine, detail) + "\n"))
	p.writer.Write([]byte(spacer + "\n"))
	p.writer.Write([]byte(fmt.Sprintf(promptLine, "Allow? (y)es / (n)o / (a)lways") + "\n"))
	p.writer.Write([]byte(border + "\n"))
}

type DefaultPrompter struct{}

func (p *DefaultPrompter) Decide(toolName string, input map[string]any, reason string) Decision {
	return Ask
}

func NewDefaultPrompter() *DefaultPrompter {
	return &DefaultPrompter{}
}

type StdinPrompter struct {
	in  *bufio.Reader
	out io.Writer
}

func NewStdinPrompter(in *bufio.Reader, out io.Writer) *StdinPrompter {
	return &StdinPrompter{
		in:  in,
		out: out,
	}
}

func (p *StdinPrompter) Decide(toolName string, input map[string]any, reason string) Decision {
	detail := p.formatDetail(toolName, input)

	p.printPrompt(toolName, detail)

	for {
		fmt.Fprint(p.out, "Allow? (y)es / (n)o / (a)lways: ")

		line, err := p.in.ReadString('\n')
		if err != nil {
			return Deny
		}

		line = strings.TrimSpace(strings.ToLower(line))

		switch line {
		case "y", "yes":
			return Allow
		case "n", "no":
			return Deny
		case "a", "always":
			return Allow
		}
	}
}

func (p *StdinPrompter) formatDetail(toolName string, input map[string]any) string {
	switch toolName {
	case "Bash":
		if cmd, ok := input["command"].(string); ok {
			return fmt.Sprintf("Command: %s", cmd)
		}
	case "Write", "Edit":
		if path, ok := input["file_path"].(string); ok {
			return fmt.Sprintf("File: %s", path)
		}
	case "Read":
		if path, ok := input["file_path"].(string); ok {
			return fmt.Sprintf("File: %s", path)
		}
	case "Glob":
		if pattern, ok := input["pattern"].(string); ok {
			return fmt.Sprintf("Pattern: %s", pattern)
		}
	case "Grep":
		if pattern, ok := input["pattern"].(string); ok {
			return fmt.Sprintf("Pattern: %s", pattern)
		}
	}
	return "No details available"
}

func (p *StdinPrompter) printPrompt(toolName string, detail string) {
	border := "┌──────────────────────────────────────────────┐"
	inner := "│  Tool: %-37s │"
	detailLine := "│  %-44s │"
	spacer := "│                                              │"
	promptLine := "│  %-44s │"

	fmt.Fprintln(p.out, border)
	fmt.Fprintln(p.out, fmt.Sprintf(inner, toolName))
	fmt.Fprintln(p.out, spacer)
	fmt.Fprintln(p.out, fmt.Sprintf(detailLine, detail))
	fmt.Fprintln(p.out, spacer)
	fmt.Fprintln(p.out, fmt.Sprintf(promptLine, "Allow? (y)es / (n)o / (a)lways"))
	fmt.Fprintln(p.out, border)
}
