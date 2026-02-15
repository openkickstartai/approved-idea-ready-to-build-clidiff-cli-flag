package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParseFlags(t *testing.T) {
	help := `Usage: myapp [OPTIONS]

Options:
  --verbose       Enable verbose output
  --output FILE   Output file path
  -h, --help      Show help
`
	flags := ParseFlags(help)
	want := map[string]bool{"--verbose": true, "--output": true, "--help": true}
	if len(flags) != len(want) {
		t.Fatalf("expected %d flags, got %d: %v", len(want), len(flags), flags)
	}
	for _, f := range flags {
		if !want[f] {
			t.Errorf("unexpected flag: %s", f)
		}
	}
}

func TestParseSubcommands(t *testing.T) {
	help := `Usage: myapp <command>

Available Commands:
  init        Initialize a project
  build       Build the project
  test        Run tests

Flags:
  --help   Show help
`
	cmds := ParseSubcommands(help)
	want := []string{"init", "build", "test"}
	if len(cmds) != len(want) {
		t.Fatalf("expected %d subcommands, got %d: %v", len(want), len(cmds), cmds)
	}
	for i, c := range cmds {
		if c != want[i] {
			t.Errorf("position %d: expected %s, got %s", i, want[i], c)
		}
	}
}

func TestDiffBreaking(t *testing.T) {
	old := Snapshot{
		Flags:       []string{"--verbose", "--output", "--format"},
		Subcommands: []string{"init", "build", "deploy"},
	}
	cur := Snapshot{
		Flags:       []string{"--verbose", "--output", "--json"},
		Subcommands: []string{"init", "build"},
	}
	r := Diff(old, cur)
	if !r.HasBreaking {
		t.Fatal("expected breaking changes")
	}
	if len(r.RemovedFlags) != 1 || r.RemovedFlags[0] != "--format" {
		t.Errorf("expected [--format] removed, got %v", r.RemovedFlags)
	}
	if len(r.RemovedCommands) != 1 || r.RemovedCommands[0] != "deploy" {
		t.Errorf("expected [deploy] removed, got %v", r.RemovedCommands)
	}
	if len(r.AddedFlags) != 1 || r.AddedFlags[0] != "--json" {
		t.Errorf("expected [--json] added, got %v", r.AddedFlags)
	}
}

func TestDiffNoBreaking(t *testing.T) {
	old := Snapshot{Flags: []string{"--verbose"}, Subcommands: []string{"init"}}
	cur := Snapshot{Flags: []string{"--verbose", "--debug"}, Subcommands: []string{"init", "build"}}
	r := Diff(old, cur)
	if r.HasBreaking {
		t.Fatal("expected no breaking changes")
	}
	if len(r.AddedFlags) != 1 || r.AddedFlags[0] != "--debug" {
		t.Errorf("expected [--debug] added, got %v", r.AddedFlags)
	}
}

func TestLoadSnapshot(t *testing.T) {
	snap := Snapshot{Command: "test", Flags: []string{"--help"}, Subcommands: []string{"run"}}
	data, _ := json.Marshal(snap)
	path := filepath.Join(t.TempDir(), "snap.json")
	os.WriteFile(path, data, 0644)
	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.Command != "test" || len(loaded.Flags) != 1 || len(loaded.Subcommands) != 1 {
		t.Errorf("snapshot mismatch: %+v", loaded)
	}
}
