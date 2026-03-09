package cli

import "snapgit/internal/git"

func runDiff() error {
	return git.Run("diff")
}