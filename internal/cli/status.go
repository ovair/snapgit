package cli

import "snapgit/internal/git"

func runStatus() error {
	return git.RunGitCommand("status")
}