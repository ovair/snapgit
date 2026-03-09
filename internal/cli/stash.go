package cli

import "snapgit/internal/git"

func runStash() error {
	return git.Run("stash")
}
