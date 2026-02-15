package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Snapshot struct {
	Command     string   `json:"command"`
	Flags       []string `json:"flags"`
	Subcommands []string `json:"subcommands"`
	CapturedAt  string   `json:"captured_at"`
}

type DiffResult struct {
	RemovedFlags    []string
	AddedFlags      []string
	RemovedCommands []string
	AddedCommands   []string
	HasBreaking     bool
}

var flagRe = regexp.MustCompile(`--([a-zA-Z][\w-]*)`)

func ParseFlags(helpText string) []string {
	matches := flagRe.FindAllString(helpText, -1)
	seen := map[string]bool{}
	var result []string
	for _, m := range matches {
		if !seen[m] {
			seen[m] = true
			result = append(result, m)
		}
	}
	return result
}

func ParseSubcommands(helpText string) []string {
	var result []string
	lines := strings.Split(helpText, "\n")
	inSection := false
	re := regexp.MustCompile(`^\s{2,}([a-zA-Z][\w-]*)\s`)
	for _, line := range lines {
		trimmed := strings.ToLower(strings.TrimSpace(line))
		if strings.HasSuffix(trimmed, "commands:") {
			inSection = true
			continue
		}
		if inSection {
			if strings.TrimSpace(line) == "" {
				inSection = false
				continue
			}
			if m := re.FindStringSubmatch(line); m != nil {
				result = append(result, m[1])
			}
		}
	}
	return result
}

func LoadSnapshot(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, err
	}
	var s Snapshot
	return s, json.Unmarshal(data, &s)
}

func Diff(old, cur Snapshot) DiffResult {
	r := DiffResult{
		RemovedFlags:    subtract(old.Flags, cur.Flags),
		AddedFlags:      subtract(cur.Flags, old.Flags),
		RemovedCommands: subtract(old.Subcommands, cur.Subcommands),
		AddedCommands:   subtract(cur.Subcommands, old.Subcommands),
	}
	r.HasBreaking = len(r.RemovedFlags) > 0 || len(r.RemovedCommands) > 0
	return r
}

func subtract(a, b []string) []string {
	set := map[string]bool{}
	for _, x := range b {
		set[x] = true
	}
	var out []string
	for _, x := range a {
		if !set[x] {
			out = append(out, x)
		}
	}
	return out
}

func PrintDiff(r DiffResult) {
	if !r.HasBreaking && len(r.AddedFlags) == 0 && len(r.AddedCommands) == 0 {
		fmt.Println("✅ No changes detected.")
		return
	}
	for _, f := range r.RemovedFlags {
		fmt.Printf("❌ BREAKING: flag removed: %s\n", f)
	}
	for _, c := range r.RemovedCommands {
		fmt.Printf("❌ BREAKING: subcommand removed: %s\n", c)
	}
	for _, f := range r.AddedFlags {
		fmt.Printf("✅ Added flag: %s\n", f)
	}
	for _, c := range r.AddedCommands {
		fmt.Printf("✅ Added subcommand: %s\n", c)
	}
	if r.HasBreaking {
		fmt.Println("\n⚠️  Breaking changes detected! Exit code 1.")
	}
}
