package cli

import (
	"fmt"
	"os"
)

func Execute() error {
	if len(os.Args) < 2 {
		printHelp()
		return nil
	}

	command := os.Args[1]

	switch command {
	case "create":
		return runCreate()
	case "get":
		return runGet()
	case "status":
		return runStatus()
	case "save":
		return runSave()
	case "diff":
		return runDiff()
	case "log":
		return runLog()
	case "branch":
		return runBranch()
	case "new":
		return runNew()
	case "go":
		return runGo()
	case "fetch":
		return runFetch()
	case "pull":
		return runPull()
	case "send":
		return runSend()
	case "help":
		printHelp()
		return nil
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func printHelp() {
	fmt.Println("SnapGit (sg)")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  create           Create a new repository")
	fmt.Println("  get <url>        Clone a repository")
	fmt.Println("  status           Show repository status")
	fmt.Println("  save <message>   Save all changes")
	fmt.Println("  diff             Show current changes")
	fmt.Println("  log              Show commit history")
	fmt.Println("  branch           List branches")
	fmt.Println("  new <branch>     Create and switch to a new branch")
	fmt.Println("  go <branch>      Switch to a branch")
	fmt.Println("  fetch            Fetch remote changes")
	fmt.Println("  pull             Pull remote changes")
	fmt.Println("  send             Push local commits")
}