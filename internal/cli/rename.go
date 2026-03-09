package cli

import (
	"fmt"
	"os"

	"snapgit/internal/git"
)

func runRename() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg rename <new-name>")
	}
	return git.Run("branch", "-m", os.Args[2])
}
