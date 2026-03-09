package cli

import "snapgit/internal/git"

func runUndo() error {
	return git.Run("reset", "--soft", "HEAD~1")
}
