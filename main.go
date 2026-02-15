package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: clidiff <snapshot|diff> ...")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "snapshot":
		runSnapshot()
	case "diff":
		runDiff()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\nUsage: clidiff <snapshot|diff>\n", os.Args[1])
		os.Exit(1)
	}
}

func runSnapshot() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: clidiff snapshot <command> [-o file.json]")
		os.Exit(1)
	}
	cmdName := os.Args[2]
	outFile := ""
	if len(os.Args) >= 5 && os.Args[3] == "-o" {
		outFile = os.Args[4]
	}
	cmd := exec.Command(cmdName, "--help")
	out, _ := cmd.CombinedOutput()
	snap := Snapshot{
		Command:     cmdName,
		Flags:       ParseFlags(string(out)),
		Subcommands: ParseSubcommands(string(out)),
		CapturedAt:  time.Now().UTC().Format(time.RFC3339),
	}
	data, _ := json.MarshalIndent(snap, "", "  ")
	if outFile != "" {
		os.WriteFile(outFile, data, 0644)
		fmt.Fprintf(os.Stderr, "Snapshot saved to %s (%d flags, %d subcommands)\n",
			outFile, len(snap.Flags), len(snap.Subcommands))
	} else {
		fmt.Println(string(data))
	}
}

func runDiff() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "Usage: clidiff diff <old.json> <new.json>")
		os.Exit(1)
	}
	old, err := LoadSnapshot(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", os.Args[2], err)
		os.Exit(1)
	}
	cur, err := LoadSnapshot(os.Args[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", os.Args[3], err)
		os.Exit(1)
	}
	result := Diff(old, cur)
	PrintDiff(result)
	if result.HasBreaking {
		os.Exit(1)
	}
}
