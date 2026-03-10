package cli

import (
	"os"

	"snapgit/internal/git"
)

func runStash() error {
	if len(os.Args) >= 3 {
		return git.Run("stash", "push", "-m", os.Args[2])
	}
	return git.Run("stash")
}
