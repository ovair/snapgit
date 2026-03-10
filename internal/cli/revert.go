package cli

import (
	"fmt"
	"os"

	"snapgit/internal/git"
)

func runRevert() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg revert <commit>")
	}
	return git.Run("revert", os.Args[2])
}
