package cli

import (
	"fmt"
	"os"
	"snapgit/internal/git"
)

func runGo() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg go <branch>")
	}
	return git.Run("switch", os.Args[2])
}
