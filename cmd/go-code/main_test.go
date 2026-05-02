package main

import (
	"bytes"
	"flag"
	"strings"
	"testing"
)

func TestRootHelpListsTaskEntrypoints(t *testing.T) {
	var out bytes.Buffer
	flags, _ := newRootFlagSet("go-code", flag.ContinueOnError, &out)

	err := flags.Parse([]string{"--help"})
	if err != flag.ErrHelp {
		t.Fatalf("Parse(--help) error = %v, want flag.ErrHelp", err)
	}

	output := out.String()
	for _, want := range []string{
		"Usage: go-code [options]",
		"go-code doctor [--offline]",
		"interactive mode",
		"setup",
		"doctor",
		"prompt mode",
		"JSON output",
		"quiet mode",
		"debug mode",
		"permission-mode",
		"version",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected %q in help output:\n%s", want, output)
		}
	}
}

func TestPermissionModeFlagValid(t *testing.T) {
	var out bytes.Buffer
	for _, mode := range []string{"read-only", "workspace-write", "danger-full-access"} {
		flags, opts := newRootFlagSet("go-code", flag.ContinueOnError, &out)
		if err := flags.Parse([]string{"--permission-mode", mode}); err != nil {
			t.Fatalf("Parse(--permission-mode %s) error = %v", mode, err)
		}
		if opts.permissionMode != mode {
			t.Errorf("expected permissionMode %q, got %q", mode, opts.permissionMode)
		}
	}
}

func TestPermissionModeFlagDefault(t *testing.T) {
	var out bytes.Buffer
	flags, opts := newRootFlagSet("go-code", flag.ContinueOnError, &out)
	if err := flags.Parse([]string{}); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if opts.permissionMode != "" {
		t.Errorf("expected default permissionMode '', got %q", opts.permissionMode)
	}
}

func TestRootVersionFlag(t *testing.T) {
	var out bytes.Buffer
	flags, opts := newRootFlagSet("go-code", flag.ContinueOnError, &out)
	if err := flags.Parse([]string{"--version"}); err != nil {
		t.Fatalf("Parse(--version) error = %v", err)
	}
	if !opts.version {
		t.Fatal("expected version option to be true")
	}
}

func TestPrintVersion(t *testing.T) {
	var out bytes.Buffer
	printVersion(&out)
	if got, want := strings.TrimSpace(out.String()), "go-code "+version; got != want {
		t.Fatalf("printVersion() = %q, want %q", got, want)
	}
}
