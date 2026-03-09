package cli

import (
	"fmt"
	"os"
	"snapgit/internal/git"
)

func runNew() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg new <branch>")
	}
	return git.Run("switch", "-c", os.Args[2])
}