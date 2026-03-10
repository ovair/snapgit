package cli

import (
	"fmt"
	"os"

	"snapgit/internal/git"
)

func runCherryPick() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg cherry-pick <commit>")
	}
	return git.Run("cherry-pick", os.Args[2])
}
