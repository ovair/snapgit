package cli

import (
	"fmt"
	"os"
	"snapgit/internal/git"
)

func runGet() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg get <url> [directory]")
	}
	args := append([]string{"clone"}, os.Args[2:]...)
	return git.Run(args...)
}
