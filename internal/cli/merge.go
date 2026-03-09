package cli

import (
	"fmt"
	"os"

	"snapgit/internal/git"
)

func runMerge() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg merge <branch>")
	}
	return git.Run("merge", os.Args[2])
}
