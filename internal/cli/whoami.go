package cli

import (
	"fmt"
	"strings"

	"snapgit/internal/git"
)

func runWhoami() error {
	name, err := git.RunOutput("config", "user.name")
	if err != nil {
		return fmt.Errorf("git user.name is not set")
	}
	email, err := git.RunOutput("config", "user.email")
	if err != nil {
		return fmt.Errorf("git user.email is not set")
	}
	fmt.Printf("%s <%s>\n", strings.TrimSpace(name), strings.TrimSpace(email))
	return nil
}
