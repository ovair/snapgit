package git

import (
	"bytes"
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

// execCmd runs a command with signal forwarding and exit code preservation.
func execCmd(cmd *exec.Cmd) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, forwardedSignals...)
	defer signal.Stop(sigCh)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", cmd.Path, err)
	}

	done := make(chan struct{})
	go func() {
		select {
		case sig := <-sigCh:
			_ = cmd.Process.Signal(sig)
		case <-done:
		}
	}()

	err := cmd.Wait()
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

// Run is the function used to execute git commands.
// It can be replaced in tests to avoid shelling out.
var Run = run

func run(args ...string) error {
	path, err := resolveGit()
	if err != nil {
		return fmt.Errorf("git is not installed or not in PATH: %w", err)
	}

	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return execCmd(cmd)
}

// RunOutput executes a git command and returns its stdout as a string.
// It can be replaced in tests.
var RunOutput = runOutput

func runOutput(args ...string) (string, error) {
	path, err := resolveGit()
	if err != nil {
		return "", fmt.Errorf("git is not installed or not in PATH: %w", err)
	}

	var stdout bytes.Buffer
	cmd := exec.Command(path, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := execCmd(cmd); err != nil {
		return "", err
	}
	return stdout.String(), nil
}

// RunExternal executes an arbitrary command (not git) with signal forwarding
// and exit code preservation. Used for tools like gh.
var RunExternal = runExternal

func runExternal(name string, args ...string) error {
	path, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("%s is not installed or not in PATH: %w", name, err)
	}

	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return execCmd(cmd)
}
