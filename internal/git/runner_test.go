package git

import (
	"errors"
	"fmt"
	"testing"
)

func TestExitError_Error(t *testing.T) {
	inner := errors.New("exit status 128")
	e := &ExitError{Code: 128, Err: inner}

	if got := e.Error(); got != "exit status 128" {
		t.Errorf("Error() = %q, want %q", got, "exit status 128")
	}
}

func TestExitError_Unwrap(t *testing.T) {
	inner := errors.New("something failed")
	e := &ExitError{Code: 1, Err: inner}

	if got := e.Unwrap(); got != inner {
		t.Errorf("Unwrap() returned different error")
	}
}

func TestExitCode_WithExitError(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{"code 0", 0},
		{"code 1", 1},
		{"code 128", 128},
		{"code 255", 255},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &ExitError{Code: tt.code, Err: errors.New("test")}
			if got := ExitCode(err); got != tt.code {
				t.Errorf("ExitCode() = %d, want %d", got, tt.code)
			}
		})
	}
}

func TestExitCode_WithWrappedExitError(t *testing.T) {
	inner := &ExitError{Code: 42, Err: errors.New("inner")}
	wrapped := fmt.Errorf("outer: %w", inner)

	if got := ExitCode(wrapped); got != 42 {
		t.Errorf("ExitCode(wrapped) = %d, want 42", got)
	}
}

func TestExitCode_WithPlainError(t *testing.T) {
	err := errors.New("plain error")
	if got := ExitCode(err); got != 1 {
		t.Errorf("ExitCode(plain) = %d, want 1", got)
	}
}

func TestExitCode_ErrorsAs(t *testing.T) {
	// Verify ExitError satisfies errors.As chain
	inner := &ExitError{Code: 5, Err: errors.New("deep")}
	wrapped := fmt.Errorf("layer1: %w", fmt.Errorf("layer2: %w", inner))

	var target *ExitError
	if !errors.As(wrapped, &target) {
		t.Fatal("errors.As failed to find ExitError in chain")
	}
	if target.Code != 5 {
		t.Errorf("target.Code = %d, want 5", target.Code)
	}
}

func TestResolveGit(t *testing.T) {
	// git should be available in the test environment
	path, err := resolveGit()
	if err != nil {
		t.Skipf("git not in PATH: %v", err)
	}
	if path == "" {
		t.Error("resolveGit() returned empty path")
	}
}

func TestRunGitCommand_Version(t *testing.T) {
	// Integration test: `git --version` should always succeed
	err := run("--version")
	if err != nil {
		t.Errorf("run(--version) failed: %v", err)
	}
}

func TestRunGitCommand_InvalidSubcommand(t *testing.T) {
	err := run("not-a-real-subcommand-xyz")
	if err == nil {
		t.Fatal("expected error for invalid git subcommand")
	}

	var exitErr *ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected ExitError, got %T: %v", err, err)
	}
	if exitErr.Code == 0 {
		t.Error("expected non-zero exit code")
	}
}

func TestRun_IsSwappable(t *testing.T) {
	// Verify that Run is a variable that can be overridden
	called := false
	original := Run
	defer func() { Run = original }()

	Run = func(args ...string) error {
		called = true
		return nil
	}

	err := Run("anything")
	if err != nil {
		t.Errorf("mock Run returned error: %v", err)
	}
	if !called {
		t.Error("mock Run was not called")
	}
}
