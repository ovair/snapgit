package cli

import "snapgit/internal/git"

func runStatus() error {
	return git.Run("status")
}