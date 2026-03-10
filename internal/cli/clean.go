package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"snapgit/internal/git"
)

func runClean() error {
	// Check for --force flag
	force := len(os.Args) >= 3 && os.Args[2] == "--force"

	if !force {
		// Dry run: show what would be deleted
		output, err := git.RunOutput("clean", "-fdn")
		if err != nil {
			return fmt.Errorf("failed to preview clean: %w", err)
		}
		output = strings.TrimSpace(output)
		if output == "" {
			fmt.Println("nothing to clean")
			return nil
		}
		fmt.Println(output)
		fmt.Print("\nRemove these files? [y/N] ")

		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("aborted")
			return nil
		}
	}

	return git.Run("clean", "-fd")
}
