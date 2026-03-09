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
	"send":   {runSend, "sg send", "Push local commits to the remote repository.\nAutomatically sets the upstream for new branches.\n\nEquivalent to: git push -u origin HEAD"},
	"undo":   {runUndo, "sg undo", "Undo the last commit but keep all changes staged.\n\nEquivalent to: git reset --soft HEAD~1"},
	"stash":  {runStash, "sg stash", "Temporarily shelve changes in your working directory.\n\nEquivalent to: git stash"},
	"pop":    {runPop, "sg pop", "Restore the most recently stashed changes.\n\nEquivalent to: git stash pop"},
	"merge":  {runMerge, "sg merge <branch>", "Merge another branch into the current branch.\n\nEquivalent to: git merge <branch>"},
	"tag":    {runTag, "sg tag [name]", "List tags or create a new tag. Without arguments, lists all tags.\n\nEquivalent to: git tag [name]"},
	"pr":     {runPR, "sg pr [\"title\"]", "Push the current branch and create a GitHub pull request.\nWithout arguments, creates a PR with title and body filled from commits.\nWith a title argument, uses that as the PR title.\n\nRequires: gh CLI (https://cli.github.com)\nEquivalent to: git push -u origin HEAD && gh pr create --fill-verbose"},
	"rename": {runRename, "sg rename <new-name>", "Rename the current branch.\n\nEquivalent to: git branch -m <new-name>"},
	"delete": {runDelete, "sg delete <branch>", "Delete a local branch. Uses safe delete — fails if the branch has unmerged changes.\nCannot delete the branch you are currently on.\n\nEquivalent to: git branch -d <branch>"},
	"ignore": {runIgnore, "sg ignore <pattern>", "Add a pattern to the .gitignore file at the repository root.\nCreates .gitignore if it does not exist.\n\nExample: sg ignore \"*.log\""},
	"whoami": {runWhoami, "sg whoami", "Show the currently configured Git user name and email.\n\nEquivalent to: git config user.name / git config user.email"},
	"remote": {runRemote, "sg remote", "Show remote repository URLs.\n\nEquivalent to: git remote -v"},
	"amend":       {runAmend, "sg amend [\"message\"]", "Amend the last commit. With a message, rewrites the commit message.\nWithout arguments, adds staged changes to the last commit keeping the message.\n\nEquivalent to: git commit --amend -m \"message\" / git commit --amend --no-edit"},
	"release":     {runRelease, "sg release <version>", "Create a version tag and push it to the remote.\nThis triggers your release pipeline (e.g. goreleaser).\n\nEquivalent to: git tag <version> && git push origin <version>"},
	"completions": {runCompletions, "sg completions <shell>", "Generate shell completion scripts.\nSupported shells: bash, zsh, fish, powershell.\n\nUsage:\n  eval \"$(sg completions bash)\"   # bash\n  eval \"$(sg completions zsh)\"    # zsh\n  sg completions fish | source    # fish\n  sg completions powershell | Out-String | Invoke-Expression  # powershell"},
}

// commandOrder defines the display order for help output.
var commandOrder = []string{
	"create", "get", "status", "add", "save", "diff", "log",
	"branch", "new", "go", "fetch", "pull", "send",
	"undo", "stash", "pop", "merge", "tag", "pr",
	"rename", "delete", "ignore", "whoami", "remote", "amend",
	"release", "completions",
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
	"pr":     "Push and create a pull request",
	"rename": "Rename the current branch",
	"delete": "Delete a local branch",
	"ignore": "Add a pattern to .gitignore",
	"whoami": "Show git user config",
	"remote": "Show remote URLs",
	"amend":       "Amend the last commit",
	"release":     "Tag and push a version",
	"completions": "Generate shell completions",
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
	fmt.Print("SnapGit (sg) — a human-friendly git CLI\n\nUsage: sg <command> [arguments]\n\nCommands:\n")

	// Extract the argument hint from the usage string (e.g. "sg get <url>" -> "<url>")
	for _, name := range commandOrder {
		cmd := commands[name]
		label := name
		if usage := cmd.usage; len(usage) > len("sg "+name) {
			label = name + " " + usage[len("sg "+name)+1:]
		}
		desc := shortDesc[name]
		fmt.Printf("  %-20s %s\n", label, desc)
	}

	fmt.Print("  help [command]       Show help (or help for a command)\n")
	fmt.Print("  version              Show version\n")
	fmt.Print("\nRun 'sg help <command>' for more details on a specific command.\n")
}

func printCommandHelp(name string) error {
	cmd, ok := commands[name]
	if !ok {
		return fmt.Errorf("unknown command: %s\nRun 'sg help' to see available commands", name)
	}
	fmt.Printf("Usage: %s\n\n%s\n", cmd.usage, cmd.help)
	return nil
}
