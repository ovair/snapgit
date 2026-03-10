package cli

import (
	"fmt"
	"os"
	"regexp"

	"snapgit/internal/git"
)

var versionRe = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

func runRelease() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg release <version>")
	}
	version := os.Args[2]
	if !versionRe.MatchString(version) {
		return fmt.Errorf("version must be in format vX.Y.Z (e.g. v1.0.0)")
	}

	if err := git.Run("tag", version); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	if err := git.Run("push", "origin", version); err != nil {
		return fmt.Errorf("failed to push tag: %w", err)
	}

	return nil
}
