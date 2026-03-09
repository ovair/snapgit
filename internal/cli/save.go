package cli

import (
	"fmt"
	"os"
	"snapgit/internal/git"
)

func runSave() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg save \"message\"")
	}

	message := os.Args[2]

	if err := git.RunGitCommand("add", "-A"); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	if err := git.RunGitCommand("commit", "-m", message); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
