package cli

import (
	"os"

	"snapgit/internal/git"
)

func runDiff() error {
	if len(os.Args) >= 3 && os.Args[2] == "staged" {
		return git.Run("diff", "--cached")
	}
	return git.Run("diff")
}
