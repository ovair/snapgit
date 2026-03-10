package cli

import (
	"fmt"
	"os"
	"strings"

	"snapgit/internal/git"
)

func runDelete() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg delete <branch>")
	}
	branch := os.Args[2]

	// Prevent deleting the current branch by checking which branch is active.
	current, err := git.RunOutput("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return fmt.Errorf("failed to determine current branch: %w", err)
	}
	current = strings.TrimSpace(current)
	if current == "" {
		return fmt.Errorf("failed to determine current branch: empty result")
	}
	if current == branch {
		return fmt.Errorf("cannot delete the current branch %q — switch to another branch first", branch)
	}

	return git.Run("branch", "-d", branch)
}
