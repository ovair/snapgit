package cli

import (
	"fmt"
	"os"

	"snapgit/internal/git"
)

func runPR() error {
	// Step 1: Push current branch to origin
	if err := git.Run("push", "-u", "origin", "HEAD"); err != nil {
		return fmt.Errorf("failed to push branch: %w", err)
	}

	// Step 2: Create PR via gh CLI
	args := []string{"pr", "create"}
	if len(os.Args) >= 3 {
		// sg pr "title" — use provided title, fill body from commits
		args = append(args, "--title", os.Args[2], "--fill-verbose")
	} else {
		// sg pr — interactive mode
		args = append(args, "--fill-verbose")
	}

	if err := git.RunExternal("gh", args...); err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	return nil
}
