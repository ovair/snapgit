package cli

import (
	"fmt"
	"os"
	"os/exec"

	"snapgit/internal/git"
)

var lookPath = exec.LookPath

func runPR() error {
	// Check that gh is installed before pushing
	if _, err := lookPath("gh"); err != nil {
		return fmt.Errorf("gh CLI is required but not installed (https://cli.github.com)")
	}

	// Step 1: Push current branch to origin
	if err := git.Run("push", "-u", "origin", "HEAD"); err != nil {
		return fmt.Errorf("failed to push branch: %w", err)
	}

	// Step 2: Create PR via gh CLI
	args := []string{"pr", "create"}
	if len(os.Args) >= 3 {
		args = append(args, "--title", os.Args[2], "--fill-verbose")
	} else {
		args = append(args, "--fill-verbose")
	}

	if err := git.RunExternal("gh", args...); err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	return nil
}
