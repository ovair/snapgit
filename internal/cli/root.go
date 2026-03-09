package cli

import (
	"fmt"
	"os"
)

// Version is set at build time via ldflags.
var Version = "dev"

// command holds a handler and its help text.
type command struct {
	handler func() error
	usage   string
	help    string
}

// commands maps command names to their definition.
var commands = map[string]command{
	"create": {runCreate, "sg create", "Initialize a new git repository in the current directory.\n\nEquivalent to: git init"},
	"get":    {runGet, "sg get <url>", "Clone a remote repository to your local machine.\n\nEquivalent to: git clone <url>"},
	"status": {runStatus, "sg status", "Show the current state of your working directory and staged changes.\n\nEquivalent to: git status"},
	"add":    {runAdd, "sg add <file|.>", "Stage specific files for the next commit. Use '.' to stage everything.\n\nEquivalent to: git add <file>"},
	"save":   {runSave, "sg save \"message\"", "Stage all changes and commit them with a message in one step.\n\nEquivalent to: git add -A && git commit -m \"message\""},
	"diff":   {runDiff, "sg diff", "Show unstaged changes in your working directory.\n\nEquivalent to: git diff"},
	"log":    {runLog, "sg log", "Show commit history as a compact, decorated graph.\n\nEquivalent to: git log --oneline --graph --decorate"},
	"branch": {runBranch, "sg branch", "List all local branches. The current branch is highlighted.\n\nEquivalent to: git branch"},
	"new":    {runNew, "sg new <branch>", "Create a new branch and switch to it immediately.\n\nEquivalent to: git switch -c <branch>"},
	"go":     {runGo, "sg go <branch>", "Switch to an existing branch.\n\nEquivalent to: git switch <branch>"},
	"fetch":  {runFetch, "sg fetch", "Download objects and refs from the remote without merging.\n\nEquivalent to: git fetch"},
	"pull":   {runPull, "sg pull", "Fetch and merge remote changes into the current branch.\n\nEquivalent to: git pull"},
	"send":   {runSend, "sg send", "Push local commits to the remote repository.\n\nEquivalent to: git push"},
	"undo":   {runUndo, "sg undo", "Undo the last commit but keep all changes staged.\n\nEquivalent to: git reset --soft HEAD~1"},
	"stash":  {runStash, "sg stash", "Temporarily shelve changes in your working directory.\n\nEquivalent to: git stash"},
	"pop":    {runPop, "sg pop", "Restore the most recently stashed changes.\n\nEquivalent to: git stash pop"},
	"merge":  {runMerge, "sg merge <branch>", "Merge another branch into the current branch.\n\nEquivalent to: git merge <branch>"},
	"tag":    {runTag, "sg tag [name]", "List tags or create a new tag. Without arguments, lists all tags.\n\nEquivalent to: git tag [name]"},
}

// commandOrder defines the display order for help output.
var commandOrder = []string{
	"create", "get", "status", "add", "save", "diff", "log",
	"branch", "new", "go", "fetch", "pull", "send",
	"undo", "stash", "pop", "merge", "tag",
}

// short descriptions for the help listing
var shortDesc = map[string]string{
	"create": "Create a new repository",
	"get":    "Clone a repository",
	"status": "Show repository status",
	"add":    "Add file(s) to the next save",
	"save":   "Save all changes",
	"diff":   "Show current changes",
	"log":    "Show commit history",
	"branch": "List branches",
	"new":    "Create and switch to a new branch",
	"go":     "Switch to a branch",
	"fetch":  "Fetch remote changes",
	"pull":   "Pull remote changes",
	"send":   "Push local commits",
	"undo":   "Undo the last commit (keep changes)",
	"stash":  "Stash working changes",
	"pop":    "Restore stashed changes",
	"merge":  "Merge a branch",
	"tag":    "List or create tags",
}

func Execute() error {
	if len(os.Args) < 2 {
		printHelp()
		return nil
	}

	name := os.Args[1]

	switch name {
	case "help", "--help", "-h":
		if len(os.Args) >= 3 {
			return printCommandHelp(os.Args[2])
		}
		printHelp()
		return nil
	case "version", "--version", "-v":
		fmt.Println("sg version " + Version)
		return nil
	}

	cmd, ok := commands[name]
	if !ok {
		return fmt.Errorf("unknown command: %s\nRun 'sg help' to see available commands", name)
	}
	return cmd.handler()
}

func printHelp() {
	fmt.Print(`SnapGit (sg) — a human-friendly git CLI

Usage: sg <command> [arguments]

Commands:
  create             Create a new repository
  get <url>          Clone a repository
  status             Show repository status
  add <file|.>       Add file(s) to the next save
  save "message"     Save all changes
  diff               Show current changes
  log                Show commit history
  branch             List branches
  new <branch>       Create and switch to a new branch
  go <branch>        Switch to a branch
  fetch              Fetch remote changes
  pull               Pull remote changes
  send               Push local commits
  undo               Undo the last commit (keep changes)
  stash              Stash working changes
  pop                Restore stashed changes
  merge <branch>     Merge a branch
  tag [name]         List or create tags
  help [command]     Show help (or help for a command)
  version            Show version

Run 'sg help <command>' for more details on a specific command.
`)
}

func printCommandHelp(name string) error {
	cmd, ok := commands[name]
	if !ok {
		return fmt.Errorf("unknown command: %s\nRun 'sg help' to see available commands", name)
	}
	fmt.Printf("Usage: %s\n\n%s\n", cmd.usage, cmd.help)
	return nil
}
