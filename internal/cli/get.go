package cli

import (
	"fmt"
	"os"
	"snapgit/internal/git"
)

func runGet() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg get <url>")
	}
	return git.Run("clone", os.Args[2])
}