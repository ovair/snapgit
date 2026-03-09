package cli

import (
	"os"

	"snapgit/internal/git"
)

func runTag() error {
	if len(os.Args) < 3 {
		// No args: list tags
		return git.Run("tag")
	}
	return git.Run("tag", os.Args[2])
}
