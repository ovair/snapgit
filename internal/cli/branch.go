package cli

import "snapgit/internal/git"

func runBranch() error {
	return git.Run("branch")
}
