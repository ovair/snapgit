package cli

import "snapgit/internal/git"

func runSend() error {
	return git.Run("push", "-u", "origin", "HEAD")
}
