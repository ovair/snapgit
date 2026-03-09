package cli

import "snapgit/internal/git"

func runDiff() error {
	return git.RunGitCommand("diff")
}