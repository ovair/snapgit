package cli

import "snapgit/internal/git"

func runPull() error {
	return git.RunGitCommand("pull")
}