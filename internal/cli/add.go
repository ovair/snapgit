package cli

import (
	"fmt"
	"os"

	"snapgit/internal/git"
)

func runAdd() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg add <file|.>")
	}
	return git.RunGitCommand("add", os.Args[2])
}
