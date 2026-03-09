package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
)

var (
	gitPath     string
	gitPathOnce sync.Once
	gitPathErr  error
)

// resolveGit finds the git binary once and caches the full path.
func resolveGit() (string, error) {
	gitPathOnce.Do(func() {
		gitPath, gitPathErr = exec.LookPath("git")
	})
	return gitPath, gitPathErr
}

// ExitError wraps an error with the git process exit code.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string { return e.Err.Error() }
func (e *ExitError) Unwrap() error { return e.Err }

// ExitCode extracts the process exit code from an error, or 1 as fallback.
func ExitCode(err error) int {
	var exitErr *ExitError
	if errors.As(err, &exitErr) {
		return exitErr.Code
	}
	return 1
}

// Run is the function used to execute git commands.
// It can be replaced in tests to avoid shelling out.
var Run = run

// run executes a git command with signal forwarding and exit code preservation.
func run(args ...string) error {
	path, err := resolveGit()
	if err != nil {
		return fmt.Errorf("git is not installed or not in PATH: %w", err)
	}

	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Forward interrupt signals to the child process.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start git: %w", err)
	}

	done := make(chan struct{})
	go func() {
		select {
		case sig := <-sigCh:
			_ = cmd.Process.Signal(sig)
		case <-done:
		}
	}()

	err = cmd.Wait()
	close(done)

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return &ExitError{Code: exitErr.ExitCode(), Err: err}
		}
		return err
	}
	return nil
}

// RunGitCommand is an alias for backward compatibility.
var RunGitCommand = Run
