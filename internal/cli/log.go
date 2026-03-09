package cli

import "snapgit/internal/git"

func runLog() error {
	return git.Run("log", "--oneline", "--graph", "--decorate")
}