package cli

import (
	"fmt"
	"os"

	"snapgit/internal/git"
)

func runRelease() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg release <version>")
	}
	version := os.Args[2]

	if err := git.Run("tag", version); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	if err := git.Run("push", "origin", version); err != nil {
		return fmt.Errorf("failed to push tag: %w", err)
	}

	return nil
}
