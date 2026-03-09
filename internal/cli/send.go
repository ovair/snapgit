package cli

import "snapgit/internal/git"

func runSend() error {
	return git.RunGitCommand("push")
}