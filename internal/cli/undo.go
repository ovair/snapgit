package cli

import (
	"fmt"
	"os"
	"strconv"

	"snapgit/internal/git"
)

func runUndo() error {
	n := "1"
	if len(os.Args) >= 3 {
		count, err := strconv.Atoi(os.Args[2])
		if err != nil || count < 1 {
			return fmt.Errorf("usage: sg undo [n] — n must be a positive integer")
		}
		n = os.Args[2]
	}
	return git.Run("reset", "--soft", "HEAD~"+n)
}
