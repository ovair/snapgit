package cli

import (
	"os"

	"snapgit/internal/git"
)

func runAmend() error {
	if len(os.Args) >= 3 {
		// sg amend "new message" — rewrite the last commit message
		return git.Run("commit", "--amend", "-m", os.Args[2])
	}
	// sg amend — add staged changes to the last commit, keep message
	return git.Run("commit", "--amend", "--no-edit")
}
