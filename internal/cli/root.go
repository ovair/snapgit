package cli

import (
	"fmt"
	"os"
)

// Version is set at build time via ldflags.
var Version = "dev"

// commands maps command names to their handler functions.
var commands = map[string]func() error{
	"create": runCreate,
	"get":    runGet,
	"status": runStatus,
	"add":    runAdd,
	"save":   runSave,
	"diff":   runDiff,
	"log":    runLog,
	"branch": runBranch,
	"new":    runNew,
	"go":     runGo,
	"fetch":  runFetch,
	"pull":   runPull,
	"send":   runSend,
}

func Execute() error {
	if len(os.Args) < 2 {
		printHelp()
		return nil
	}

	command := os.Args[1]

	if command == "help" {
		printHelp()
		return nil
	}

	if command == "version" || command == "--version" || command == "-v" {
		fmt.Println("sg version " + Version)
		return nil
	}

	fn, ok := commands[command]
	if !ok {
		return fmt.Errorf("unknown command: %s\nRun 'sg help' to see available commands", command)
	}
	return fn()
}

func printHelp() {
	fmt.Println("SnapGit (sg)")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  create           Create a new repository")
	fmt.Println("  get <url>        Clone a repository")
	fmt.Println("  status           Show repository status")
	fmt.Println("  add <file|.>     Add file(s) to the next save")
	fmt.Println("  save <message>   Save all changes")
	fmt.Println("  diff             Show current changes")
	fmt.Println("  log              Show commit history")
	fmt.Println("  branch           List branches")
	fmt.Println("  new <branch>     Create and switch to a new branch")
	fmt.Println("  go <branch>      Switch to a branch")
	fmt.Println("  fetch            Fetch remote changes")
	fmt.Println("  pull             Pull remote changes")
	fmt.Println("  send             Push local commits")
	fmt.Println("  help             Show this help message")
}
