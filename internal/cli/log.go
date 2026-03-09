package cli

import "snapgit/internal/git"

func runLog() error {
	return git.RunGitCommand("log", "--oneline", "--graph", "--decorate")
}