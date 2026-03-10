package cli

import "snapgit/internal/git"

func runCreate() error {
	return git.Run("init")
}
