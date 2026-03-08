package main

import (
	"strings"
	"testing"
)

func TestInferOwner(t *testing.T) {
	tests := []struct {
		name       string
		modulePath string
		want       string
	}{
		{"github path", "github.com/user/repo", "user"},
		{"short path", "example.com/pkg", "pkg"},
		{"single segment", "mypackage", ""},
		{"empty", "", ""},
		{"deep path", "github.com/org/repo/sub", "org"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inferOwner(tt.modulePath); got != tt.want {
				t.Errorf("inferOwner(%q) = %q, want %q", tt.modulePath, got, tt.want)
			}
		})
	}
}

func TestNewRootCmd(t *testing.T) {
	cmd := newRootCmd()

	if cmd.Use != "goseed" {
		t.Errorf("Use = %q, want %q", cmd.Use, "goseed")
	}

	flags := []string{
		"name", "module", "type", "go-version", "license", "build-tool",
		"http-framework", "linter", "goreleaser", "ci", "docker", "env-example",
		"dependabot", "claude", "claude-ci", "no-interactive", "output-dir",
	}
	for _, name := range flags {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag %q to exist", name)
		}
	}

	var foundVersion bool
	for _, sub := range cmd.Commands() {
		if sub.Use == "version" {
			foundVersion = true
			break
		}
	}
	if !foundVersion {
		t.Error("expected version subcommand")
	}
}

func TestNonInteractiveValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			"missing name",
			[]string{"--no-interactive", "--module", "github.com/x/y", "--type", "cli"},
			"--name is required",
		},
		{
			"missing module",
			[]string{"--no-interactive", "--name", "foo", "--type", "cli"},
			"--module is required",
		},
		{
			"missing type",
			[]string{"--no-interactive", "--name", "foo", "--module", "github.com/x/y"},
			"--type is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCmd()
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want containing %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestNonInteractiveGenerate(t *testing.T) {
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{
		"--no-interactive",
		"--name", "testapp",
		"--module", "github.com/test/testapp",
		"--type", "cli",
		"--output-dir", dir,
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}
