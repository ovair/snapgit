package cli

import (
	"fmt"
	"os"
	"strings"

	"snapgit/internal/git"
)

func runIgnore() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg ignore <pattern>")
	}
	pattern := os.Args[2]

	// Find the repo root so we always write to the top-level .gitignore.
	root, err := git.RunOutput("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("not a git repository")
	}
	root = strings.TrimSpace(root)

	path := root + "/.gitignore"

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore: %w", err)
	}
	defer f.Close()

	if _, err := fmt.Fprintln(f, pattern); err != nil {
		return fmt.Errorf("failed to write to .gitignore: %w", err)
	}

	fmt.Printf("added %q to .gitignore\n", pattern)
	return nil
}
