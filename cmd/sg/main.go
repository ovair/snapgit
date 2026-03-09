package main

import (
	"fmt"
	"os"
	"snapgit/internal/cli"
	"snapgit/internal/git"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(git.ExitCode(err))
	}
}
