package cli

import "snapgit/internal/git"

func runRemote() error {
	return git.Run("remote", "-v")
}
