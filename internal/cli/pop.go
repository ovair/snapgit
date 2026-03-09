package cli

import "snapgit/internal/git"

func runPop() error {
	return git.Run("stash", "pop")
}
