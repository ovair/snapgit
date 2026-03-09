package cli

import "snapgit/internal/git"

func runFetch() error {
	return git.Run("fetch")
}