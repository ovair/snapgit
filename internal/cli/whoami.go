package cli

import (
	"fmt"
	"strings"

	"snapgit/internal/git"
)

func runWhoami() error {
	name, nameErr := git.RunOutput("config", "user.name")
	email, emailErr := git.RunOutput("config", "user.email")

	var missing []string
	if nameErr != nil {
		missing = append(missing, "user.name")
	}
	if emailErr != nil {
		missing = append(missing, "user.email")
	}
	if len(missing) > 0 {
		return fmt.Errorf("git config %s not set — run: git config --global %s <value>", strings.Join(missing, " and "), missing[0])
	}

	fmt.Printf("%s <%s>\n", strings.TrimSpace(name), strings.TrimSpace(email))
	return nil
}
