package cli

import (
	"fmt"
	"os"
	"strconv"

	"snapgit/internal/git"
)

func runLog() error {
	args := []string{"log", "--oneline", "--graph", "--decorate"}
	if len(os.Args) >= 3 {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil || n < 1 {
			return fmt.Errorf("usage: sg log [n] — n must be a positive integer")
		}
		args = append(args, fmt.Sprintf("-%d", n))
	}
	return git.Run(args...)
}
