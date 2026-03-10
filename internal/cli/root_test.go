package cli

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"snapgit/internal/git"
)

// mockRun records git commands without executing them.
type mockRun struct {
	calls [][]string
	err   error // if set, all calls return this error
}

func (m *mockRun) run(args ...string) error {
	m.calls = append(m.calls, args)
	return m.err
}

// withArgs sets os.Args for the duration of a test and restores it after.
func withArgs(t *testing.T, args []string) {
	t.Helper()
	orig := os.Args
	t.Cleanup(func() { os.Args = orig })
	os.Args = args
}

// withMockGit swaps git.Run for a mock and restores it after.
func withMockGit(t *testing.T) *mockRun {
	t.Helper()
	m := &mockRun{}
	orig := git.Run
	t.Cleanup(func() { git.Run = orig })
	git.Run = m.run
	return m
}

// mockExternal records external command calls.
type mockExternal struct {
	calls []struct {
		name string
		args []string
	}
	err error
}

func (m *mockExternal) run(name string, args ...string) error {
	m.calls = append(m.calls, struct {
		name string
		args []string
	}{name, args})
	return m.err
}

// mockOutput records git output calls and returns configurable output.
type mockOutput struct {
	calls  [][]string
	output string
	err    error
}

func (m *mockOutput) run(args ...string) (string, error) {
	m.calls = append(m.calls, args)
	return m.output, m.err
}

// withMockGitOutput swaps git.RunOutput for a mock and restores it after.
func withMockGitOutput(t *testing.T, output string) *mockOutput {
	t.Helper()
	m := &mockOutput{output: output}
	orig := git.RunOutput
	t.Cleanup(func() { git.RunOutput = orig })
	git.RunOutput = m.run
	return m
}

// withMockExternal swaps git.RunExternal for a mock and restores it after.
func withMockExternal(t *testing.T) *mockExternal {
	t.Helper()
	m := &mockExternal{}
	orig := git.RunExternal
	t.Cleanup(func() { git.RunExternal = orig })
	git.RunExternal = m.run
	return m
}

// captureStdout captures stdout during f() and returns it as a string.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	f()
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

// --- Execute dispatch tests ---

func TestExecute_NoArgs(t *testing.T) {
	withArgs(t, []string{"sg"})
	out := captureStdout(t, func() {
		err := Execute()
		if err != nil {
			t.Errorf("Execute() returned error: %v", err)
		}
	})
	if !strings.Contains(out, "SnapGit") {
		t.Errorf("expected help output, got: %s", out)
	}
}

func TestExecute_Help(t *testing.T) {
	for _, arg := range []string{"help", "--help", "-h"} {
		t.Run(arg, func(t *testing.T) {
			withArgs(t, []string{"sg", arg})
			out := captureStdout(t, func() {
				err := Execute()
				if err != nil {
					t.Errorf("Execute(%s) returned error: %v", arg, err)
				}
			})
			if !strings.Contains(out, "Commands:") {
				t.Errorf("help output missing 'Commands:': %s", out)
			}
		})
	}
}

func TestExecute_Version(t *testing.T) {
	for _, arg := range []string{"version", "--version", "-v"} {
		t.Run(arg, func(t *testing.T) {
			withArgs(t, []string{"sg", arg})
			out := captureStdout(t, func() {
				err := Execute()
				if err != nil {
					t.Errorf("Execute(%s) returned error: %v", arg, err)
				}
			})
			if !strings.Contains(out, "sg version") {
				t.Errorf("version output missing 'sg version': %s", out)
			}
		})
	}
}

func TestExecute_UnknownCommand(t *testing.T) {
	withArgs(t, []string{"sg", "notacommand"})
	err := Execute()
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExecute_PerCommandHelp(t *testing.T) {
	withArgs(t, []string{"sg", "help", "save"})
	out := captureStdout(t, func() {
		err := Execute()
		if err != nil {
			t.Errorf("Execute(help save) returned error: %v", err)
		}
	})
	if !strings.Contains(out, "sg save") {
		t.Errorf("per-command help missing usage: %s", out)
	}
	if !strings.Contains(out, "Equivalent to:") {
		t.Errorf("per-command help missing git equivalent: %s", out)
	}
}

func TestExecute_PerCommandHelp_Unknown(t *testing.T) {
	withArgs(t, []string{"sg", "help", "notreal"})
	err := Execute()
	if err == nil {
		t.Fatal("expected error for help on unknown command")
	}
}

// --- Command handler tests ---

func TestRunCreate(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "create"})

	if err := runCreate(); err != nil {
		t.Fatal(err)
	}
	if len(m.calls) != 1 || m.calls[0][0] != "init" {
		t.Errorf("expected [init], got %v", m.calls)
	}
}

func TestRunGet(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "get", "https://github.com/test/repo"})

	if err := runGet(); err != nil {
		t.Fatal(err)
	}
	if len(m.calls) != 1 || m.calls[0][0] != "clone" || m.calls[0][1] != "https://github.com/test/repo" {
		t.Errorf("expected [clone url], got %v", m.calls)
	}
}

func TestRunGet_WithDir(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "get", "https://github.com/test/repo", "mydir"})

	if err := runGet(); err != nil {
		t.Fatal(err)
	}
	args := m.calls[0]
	if args[0] != "clone" || args[1] != "https://github.com/test/repo" || args[2] != "mydir" {
		t.Errorf("expected [clone url mydir], got %v", args)
	}
}

func TestRunGet_NoArgs(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "get"})

	err := runGet()
	if err == nil {
		t.Fatal("expected error when no URL provided")
	}
}

func TestRunStatus(t *testing.T) {
	m := withMockGit(t)
	if err := runStatus(); err != nil {
		t.Fatal(err)
	}
	if len(m.calls) != 1 || m.calls[0][0] != "status" {
		t.Errorf("expected [status], got %v", m.calls)
	}
}

func TestRunAdd(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "add", "."})

	if err := runAdd(); err != nil {
		t.Fatal(err)
	}
	if len(m.calls) != 1 || m.calls[0][0] != "add" || m.calls[0][1] != "." {
		t.Errorf("expected [add .], got %v", m.calls)
	}
}

func TestRunAdd_MultipleFiles(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "add", "file1.go", "file2.go", "file3.go"})

	if err := runAdd(); err != nil {
		t.Fatal(err)
	}
	args := m.calls[0]
	if args[0] != "add" || args[1] != "file1.go" || args[2] != "file2.go" || args[3] != "file3.go" {
		t.Errorf("expected [add file1.go file2.go file3.go], got %v", args)
	}
}

func TestRunAdd_NoArgs(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "add"})

	err := runAdd()
	if err == nil {
		t.Fatal("expected error when no file provided")
	}
}

func TestRunSave(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "save", "my commit message"})

	if err := runSave(); err != nil {
		t.Fatal(err)
	}
	if len(m.calls) != 2 {
		t.Fatalf("expected 2 git calls, got %d", len(m.calls))
	}
	// First call: git add -A
	if m.calls[0][0] != "add" || m.calls[0][1] != "-A" {
		t.Errorf("first call: expected [add -A], got %v", m.calls[0])
	}
	// Second call: git commit -m "message"
	if m.calls[1][0] != "commit" || m.calls[1][1] != "-m" || m.calls[1][2] != "my commit message" {
		t.Errorf("second call: expected [commit -m msg], got %v", m.calls[1])
	}
}

func TestRunSave_NoArgs(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "save"})

	err := runSave()
	if err == nil {
		t.Fatal("expected error when no message provided")
	}
}

func TestRunSave_StageFailure(t *testing.T) {
	m := withMockGit(t)
	m.err = fmt.Errorf("staging failed")
	withArgs(t, []string{"sg", "save", "msg"})

	err := runSave()
	if err == nil {
		t.Fatal("expected error when staging fails")
	}
	if !strings.Contains(err.Error(), "failed to stage") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunDiff(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "diff"})
	if err := runDiff(); err != nil {
		t.Fatal(err)
	}
	if len(m.calls[0]) != 1 || m.calls[0][0] != "diff" {
		t.Errorf("expected [diff], got %v", m.calls)
	}
}

func TestRunDiff_Staged(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "diff", "staged"})
	if err := runDiff(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "diff" || m.calls[0][1] != "--cached" {
		t.Errorf("expected [diff --cached], got %v", m.calls[0])
	}
}

func TestRunLog(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "log"})
	if err := runLog(); err != nil {
		t.Fatal(err)
	}
	expected := []string{"log", "--oneline", "--graph", "--decorate"}
	if len(m.calls[0]) != 4 {
		t.Fatalf("expected 4 args, got %v", m.calls[0])
	}
	for i, arg := range expected {
		if m.calls[0][i] != arg {
			t.Errorf("arg[%d] = %q, want %q", i, m.calls[0][i], arg)
		}
	}
}

func TestRunLog_WithN(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "log", "10"})
	if err := runLog(); err != nil {
		t.Fatal(err)
	}
	if len(m.calls[0]) != 5 {
		t.Fatalf("expected 5 args, got %v", m.calls[0])
	}
	if m.calls[0][4] != "-10" {
		t.Errorf("expected -10, got %s", m.calls[0][4])
	}
}

func TestRunLog_InvalidN(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "log", "abc"})
	err := runLog()
	if err == nil {
		t.Fatal("expected error for invalid n")
	}
}

func TestRunBranch(t *testing.T) {
	m := withMockGit(t)
	if err := runBranch(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "branch" {
		t.Errorf("expected [branch], got %v", m.calls)
	}
}

func TestRunNew(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "new", "feature-x"})

	if err := runNew(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "switch" || m.calls[0][1] != "-c" || m.calls[0][2] != "feature-x" {
		t.Errorf("expected [switch -c feature-x], got %v", m.calls[0])
	}
}

func TestRunNew_NoArgs(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "new"})

	err := runNew()
	if err == nil {
		t.Fatal("expected error when no branch name provided")
	}
}

func TestRunGo(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "go", "main"})

	if err := runGo(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "switch" || m.calls[0][1] != "main" {
		t.Errorf("expected [switch main], got %v", m.calls[0])
	}
}

func TestRunGo_NoArgs(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "go"})

	err := runGo()
	if err == nil {
		t.Fatal("expected error when no branch name provided")
	}
}

func TestRunFetch(t *testing.T) {
	m := withMockGit(t)
	if err := runFetch(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "fetch" {
		t.Errorf("expected [fetch], got %v", m.calls)
	}
}

func TestRunPull(t *testing.T) {
	m := withMockGit(t)
	if err := runPull(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "pull" {
		t.Errorf("expected [pull], got %v", m.calls)
	}
}

func TestRunSend(t *testing.T) {
	m := withMockGit(t)
	if err := runSend(); err != nil {
		t.Fatal(err)
	}
	args := m.calls[0]
	if args[0] != "push" || args[1] != "-u" || args[2] != "origin" || args[3] != "HEAD" {
		t.Errorf("expected [push -u origin HEAD], got %v", args)
	}
}

func TestRunUndo(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "undo"})
	if err := runUndo(); err != nil {
		t.Fatal(err)
	}
	expected := []string{"reset", "--soft", "HEAD~1"}
	for i, arg := range expected {
		if m.calls[0][i] != arg {
			t.Errorf("arg[%d] = %q, want %q", i, m.calls[0][i], arg)
		}
	}
}

func TestRunUndo_WithN(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "undo", "3"})
	if err := runUndo(); err != nil {
		t.Fatal(err)
	}
	expected := []string{"reset", "--soft", "HEAD~3"}
	for i, arg := range expected {
		if m.calls[0][i] != arg {
			t.Errorf("arg[%d] = %q, want %q", i, m.calls[0][i], arg)
		}
	}
}

func TestRunUndo_InvalidN(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "undo", "abc"})
	err := runUndo()
	if err == nil {
		t.Fatal("expected error for invalid n")
	}
}

func TestRunUndo_ZeroN(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "undo", "0"})
	err := runUndo()
	if err == nil {
		t.Fatal("expected error for zero n")
	}
}

func TestRunStash(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "stash"})
	if err := runStash(); err != nil {
		t.Fatal(err)
	}
	if len(m.calls[0]) != 1 || m.calls[0][0] != "stash" {
		t.Errorf("expected [stash], got %v", m.calls)
	}
}

func TestRunStash_WithMessage(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "stash", "work in progress"})
	if err := runStash(); err != nil {
		t.Fatal(err)
	}
	args := m.calls[0]
	if args[0] != "stash" || args[1] != "push" || args[2] != "-m" || args[3] != "work in progress" {
		t.Errorf("expected [stash push -m 'work in progress'], got %v", args)
	}
}

func TestRunPop(t *testing.T) {
	m := withMockGit(t)
	if err := runPop(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "stash" || m.calls[0][1] != "pop" {
		t.Errorf("expected [stash pop], got %v", m.calls)
	}
}

func TestRunMerge(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "merge", "feature"})

	if err := runMerge(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "merge" || m.calls[0][1] != "feature" {
		t.Errorf("expected [merge feature], got %v", m.calls[0])
	}
}

func TestRunMerge_NoArgs(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "merge"})

	err := runMerge()
	if err == nil {
		t.Fatal("expected error when no branch provided")
	}
}

func TestRunTag_List(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "tag"})

	if err := runTag(); err != nil {
		t.Fatal(err)
	}
	if len(m.calls[0]) != 1 || m.calls[0][0] != "tag" {
		t.Errorf("expected [tag], got %v", m.calls[0])
	}
}

func TestRunTag_Create(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "tag", "v1.0.0"})

	if err := runTag(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "tag" || m.calls[0][1] != "v1.0.0" {
		t.Errorf("expected [tag v1.0.0], got %v", m.calls[0])
	}
}

func withMockLookPath(t *testing.T) {
	t.Helper()
	orig := lookPath
	t.Cleanup(func() { lookPath = orig })
	lookPath = func(file string) (string, error) { return "/usr/bin/" + file, nil }
}

func TestRunPR_NoTitle(t *testing.T) {
	m := withMockGit(t)
	ext := withMockExternal(t)
	withMockLookPath(t)
	withArgs(t, []string{"sg", "pr"})

	if err := runPR(); err != nil {
		t.Fatal(err)
	}

	// Should push first
	if len(m.calls) != 1 {
		t.Fatalf("expected 1 git call, got %d", len(m.calls))
	}
	push := m.calls[0]
	if push[0] != "push" || push[1] != "-u" || push[2] != "origin" || push[3] != "HEAD" {
		t.Errorf("expected [push -u origin HEAD], got %v", push)
	}

	// Should create PR with --fill-verbose
	if len(ext.calls) != 1 {
		t.Fatalf("expected 1 external call, got %d", len(ext.calls))
	}
	if ext.calls[0].name != "gh" {
		t.Errorf("expected gh, got %s", ext.calls[0].name)
	}
	ghArgs := ext.calls[0].args
	if ghArgs[0] != "pr" || ghArgs[1] != "create" || ghArgs[2] != "--fill-verbose" {
		t.Errorf("expected [pr create --fill-verbose], got %v", ghArgs)
	}
}

func TestRunPR_WithTitle(t *testing.T) {
	m := withMockGit(t)
	ext := withMockExternal(t)
	withMockLookPath(t)
	withArgs(t, []string{"sg", "pr", "Add new feature"})

	if err := runPR(); err != nil {
		t.Fatal(err)
	}

	// Should push
	if len(m.calls) != 1 || m.calls[0][0] != "push" {
		t.Errorf("expected push, got %v", m.calls)
	}

	// Should create PR with --title
	ghArgs := ext.calls[0].args
	if ghArgs[0] != "pr" || ghArgs[1] != "create" || ghArgs[2] != "--title" || ghArgs[3] != "Add new feature" || ghArgs[4] != "--fill-verbose" {
		t.Errorf("expected [pr create --title 'Add new feature' --fill-verbose], got %v", ghArgs)
	}
}

func TestRunPR_GhNotInstalled(t *testing.T) {
	withMockGit(t)
	withMockExternal(t)
	withArgs(t, []string{"sg", "pr"})
	orig := lookPath
	t.Cleanup(func() { lookPath = orig })
	lookPath = func(file string) (string, error) { return "", fmt.Errorf("not found") }

	err := runPR()
	if err == nil {
		t.Fatal("expected error when gh not installed")
	}
	if !strings.Contains(err.Error(), "gh CLI is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunPR_PushFails(t *testing.T) {
	m := withMockGit(t)
	m.err = fmt.Errorf("push failed")
	withMockExternal(t)
	withMockLookPath(t)
	withArgs(t, []string{"sg", "pr"})

	err := runPR()
	if err == nil {
		t.Fatal("expected error when push fails")
	}
	if !strings.Contains(err.Error(), "failed to push") {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- rename tests ---

func TestRunRename(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "rename", "new-name"})

	if err := runRename(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "branch" || m.calls[0][1] != "-m" || m.calls[0][2] != "new-name" {
		t.Errorf("expected [branch -m new-name], got %v", m.calls[0])
	}
}

func TestRunRename_NoArgs(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "rename"})

	err := runRename()
	if err == nil {
		t.Fatal("expected error when no name provided")
	}
}

// --- delete tests ---

func TestRunDelete(t *testing.T) {
	m := withMockGit(t)
	withMockGitOutput(t, "main\n")
	withArgs(t, []string{"sg", "delete", "feature"})

	if err := runDelete(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "branch" || m.calls[0][1] != "-d" || m.calls[0][2] != "feature" {
		t.Errorf("expected [branch -d feature], got %v", m.calls[0])
	}
}

func TestRunDelete_NoArgs(t *testing.T) {
	withMockGit(t)
	withMockGitOutput(t, "main\n")
	withArgs(t, []string{"sg", "delete"})

	err := runDelete()
	if err == nil {
		t.Fatal("expected error when no branch provided")
	}
}

func TestRunDelete_CurrentBranch(t *testing.T) {
	withMockGit(t)
	withMockGitOutput(t, "feature\n")
	withArgs(t, []string{"sg", "delete", "feature"})

	err := runDelete()
	if err == nil {
		t.Fatal("expected error when deleting current branch")
	}
	if !strings.Contains(err.Error(), "cannot delete the current branch") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunDelete_EmptyBranch(t *testing.T) {
	withMockGit(t)
	withMockGitOutput(t, "")
	withArgs(t, []string{"sg", "delete", "feature"})

	err := runDelete()
	if err == nil {
		t.Fatal("expected error when branch output is empty")
	}
	if !strings.Contains(err.Error(), "empty result") {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- ignore tests ---

func TestRunIgnore(t *testing.T) {
	dir := t.TempDir()
	withMockGit(t)
	withMockGitOutput(t, dir+"\n")
	withArgs(t, []string{"sg", "ignore", "*.log"})

	out := captureStdout(t, func() {
		if err := runIgnore(); err != nil {
			t.Fatal(err)
		}
	})

	if !strings.Contains(out, "*.log") {
		t.Errorf("expected confirmation output, got: %s", out)
	}

	data, err := os.ReadFile(dir + "/.gitignore")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "*.log") {
		t.Errorf("expected *.log in .gitignore, got: %s", string(data))
	}
}

func TestRunIgnore_NoArgs(t *testing.T) {
	withMockGit(t)
	withMockGitOutput(t, "/tmp\n")
	withArgs(t, []string{"sg", "ignore"})

	err := runIgnore()
	if err == nil {
		t.Fatal("expected error when no pattern provided")
	}
}

// --- whoami tests ---

func TestRunWhoami(t *testing.T) {
	withMockGit(t)
	mo := withMockGitOutput(t, "")
	// RunOutput is called twice: once for name, once for email
	// Override with a func that returns different values per call
	callCount := 0
	git.RunOutput = func(args ...string) (string, error) {
		mo.calls = append(mo.calls, args)
		callCount++
		if callCount == 1 {
			return "Test User\n", nil
		}
		return "test@example.com\n", nil
	}
	withArgs(t, []string{"sg", "whoami"})

	out := captureStdout(t, func() {
		if err := runWhoami(); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "Test User") || !strings.Contains(out, "test@example.com") {
		t.Errorf("expected name and email, got: %s", out)
	}
}

func TestRunWhoami_BothMissing(t *testing.T) {
	withMockGit(t)
	withMockGitOutput(t, "")
	git.RunOutput = func(args ...string) (string, error) {
		return "", fmt.Errorf("not set")
	}
	withArgs(t, []string{"sg", "whoami"})

	err := runWhoami()
	if err == nil {
		t.Fatal("expected error when both configs missing")
	}
	if !strings.Contains(err.Error(), "user.name") || !strings.Contains(err.Error(), "user.email") {
		t.Errorf("expected both missing fields reported, got: %v", err)
	}
}

func TestRunWhoami_EmailMissing(t *testing.T) {
	withMockGit(t)
	withMockGitOutput(t, "")
	callCount := 0
	git.RunOutput = func(args ...string) (string, error) {
		callCount++
		if callCount == 1 {
			return "Test User\n", nil
		}
		return "", fmt.Errorf("not set")
	}
	withArgs(t, []string{"sg", "whoami"})

	err := runWhoami()
	if err == nil {
		t.Fatal("expected error when email missing")
	}
	if !strings.Contains(err.Error(), "user.email") {
		t.Errorf("expected user.email mentioned, got: %v", err)
	}
}

// --- remote tests ---

func TestRunRemote(t *testing.T) {
	m := withMockGit(t)
	if err := runRemote(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "remote" || m.calls[0][1] != "-v" {
		t.Errorf("expected [remote -v], got %v", m.calls[0])
	}
}

// --- amend tests ---

func TestRunAmend_WithMessage(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "amend", "fixed typo"})

	if err := runAmend(); err != nil {
		t.Fatal(err)
	}
	args := m.calls[0]
	if args[0] != "commit" || args[1] != "--amend" || args[2] != "-m" || args[3] != "fixed typo" {
		t.Errorf("expected [commit --amend -m 'fixed typo'], got %v", args)
	}
}

func TestRunAmend_NoMessage(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "amend"})

	if err := runAmend(); err != nil {
		t.Fatal(err)
	}
	args := m.calls[0]
	if args[0] != "commit" || args[1] != "--amend" || args[2] != "--no-edit" {
		t.Errorf("expected [commit --amend --no-edit], got %v", args)
	}
}

// --- release tests ---

func TestRunRelease(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "release", "v0.3.0"})

	if err := runRelease(); err != nil {
		t.Fatal(err)
	}
	if len(m.calls) != 2 {
		t.Fatalf("expected 2 git calls, got %d", len(m.calls))
	}
	// First: git tag v0.3.0
	if m.calls[0][0] != "tag" || m.calls[0][1] != "v0.3.0" {
		t.Errorf("expected [tag v0.3.0], got %v", m.calls[0])
	}
	// Second: git push origin v0.3.0
	if m.calls[1][0] != "push" || m.calls[1][1] != "origin" || m.calls[1][2] != "v0.3.0" {
		t.Errorf("expected [push origin v0.3.0], got %v", m.calls[1])
	}
}

func TestRunRelease_NoArgs(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "release"})

	err := runRelease()
	if err == nil {
		t.Fatal("expected error when no version provided")
	}
}

func TestRunRelease_InvalidVersion(t *testing.T) {
	withMockGit(t)
	for _, v := range []string{"oops", "1.0.0", "vx.y.z", "v1.0"} {
		t.Run(v, func(t *testing.T) {
			withArgs(t, []string{"sg", "release", v})
			err := runRelease()
			if err == nil {
				t.Fatalf("expected error for invalid version %q", v)
			}
			if !strings.Contains(err.Error(), "vX.Y.Z") {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRunRelease_TagFails(t *testing.T) {
	m := withMockGit(t)
	m.err = fmt.Errorf("tag exists")
	withArgs(t, []string{"sg", "release", "v0.3.0"})

	err := runRelease()
	if err == nil {
		t.Fatal("expected error when tag fails")
	}
	if !strings.Contains(err.Error(), "failed to create tag") {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- completions tests ---

func TestRunCompletions_Bash(t *testing.T) {
	withArgs(t, []string{"sg", "completions", "bash"})
	out := captureStdout(t, func() {
		if err := runCompletions(); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "complete") || !strings.Contains(out, "_sg_completions") {
		t.Errorf("bash completion missing expected content: %s", out)
	}
}

func TestRunCompletions_Zsh(t *testing.T) {
	withArgs(t, []string{"sg", "completions", "zsh"})
	out := captureStdout(t, func() {
		if err := runCompletions(); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "#compdef sg") {
		t.Errorf("zsh completion missing #compdef: %s", out)
	}
}

func TestRunCompletions_Fish(t *testing.T) {
	withArgs(t, []string{"sg", "completions", "fish"})
	out := captureStdout(t, func() {
		if err := runCompletions(); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "complete -c sg") {
		t.Errorf("fish completion missing expected content: %s", out)
	}
}

func TestRunCompletions_Powershell(t *testing.T) {
	withArgs(t, []string{"sg", "completions", "powershell"})
	out := captureStdout(t, func() {
		if err := runCompletions(); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "Register-ArgumentCompleter") {
		t.Errorf("powershell completion missing expected content: %s", out)
	}
}

func TestRunCompletions_UnsupportedShell(t *testing.T) {
	withArgs(t, []string{"sg", "completions", "nushell"})
	err := runCompletions()
	if err == nil {
		t.Fatal("expected error for unsupported shell")
	}
	if !strings.Contains(err.Error(), "unsupported shell") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunCompletions_NoArgs(t *testing.T) {
	withArgs(t, []string{"sg", "completions"})
	err := runCompletions()
	if err == nil {
		t.Fatal("expected error when no shell provided")
	}
}

// --- clean tests ---

func TestRunClean_Force(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "clean", "--force"})

	if err := runClean(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "clean" || m.calls[0][1] != "-fd" {
		t.Errorf("expected [clean -fd], got %v", m.calls[0])
	}
}

func TestRunClean_NothingToClean(t *testing.T) {
	withMockGit(t)
	withMockGitOutput(t, "")
	withArgs(t, []string{"sg", "clean"})

	out := captureStdout(t, func() {
		if err := runClean(); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "nothing to clean") {
		t.Errorf("expected 'nothing to clean', got: %s", out)
	}
}

// --- revert tests ---

func TestRunRevert(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "revert", "abc123"})

	if err := runRevert(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "revert" || m.calls[0][1] != "abc123" {
		t.Errorf("expected [revert abc123], got %v", m.calls[0])
	}
}

func TestRunRevert_NoArgs(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "revert"})

	err := runRevert()
	if err == nil {
		t.Fatal("expected error when no commit provided")
	}
}

// --- cherry-pick tests ---

func TestRunCherryPick(t *testing.T) {
	m := withMockGit(t)
	withArgs(t, []string{"sg", "cherry-pick", "abc123"})

	if err := runCherryPick(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "cherry-pick" || m.calls[0][1] != "abc123" {
		t.Errorf("expected [cherry-pick abc123], got %v", m.calls[0])
	}
}

func TestRunCherryPick_NoArgs(t *testing.T) {
	withMockGit(t)
	withArgs(t, []string{"sg", "cherry-pick"})

	err := runCherryPick()
	if err == nil {
		t.Fatal("expected error when no commit provided")
	}
}

// --- ignore Windows newline test ---

func TestRunIgnore_WindowsNewline(t *testing.T) {
	dir := t.TempDir()
	withMockGit(t)
	withMockGitOutput(t, dir+"\r\n")
	withArgs(t, []string{"sg", "ignore", "*.tmp"})

	out := captureStdout(t, func() {
		if err := runIgnore(); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "*.tmp") {
		t.Errorf("expected confirmation, got: %s", out)
	}
}

// --- All commands registered ---

func TestAllCommandsRegistered(t *testing.T) {
	expected := []string{
		"create", "get", "status", "add", "save", "diff", "log",
		"branch", "new", "go", "fetch", "pull", "send",
		"undo", "stash", "pop", "merge", "tag", "pr",
		"rename", "delete", "ignore", "whoami", "remote", "amend",
		"clean", "revert", "cherry-pick",
		"release", "completions",
	}
	for _, name := range expected {
		if _, ok := commands[name]; !ok {
			t.Errorf("command %q not registered", name)
		}
	}
}

func TestCommandOrderMatchesCommands(t *testing.T) {
	for _, name := range commandOrder {
		if _, ok := commands[name]; !ok {
			t.Errorf("commandOrder contains %q which is not in commands map", name)
		}
	}
	if len(commandOrder) != len(commands) {
		t.Errorf("commandOrder has %d entries, commands has %d", len(commandOrder), len(commands))
	}
}

func TestShortDescMatchesCommands(t *testing.T) {
	for name := range commands {
		if _, ok := shortDesc[name]; !ok {
			t.Errorf("command %q is missing from shortDesc map", name)
		}
	}
	for name := range shortDesc {
		if _, ok := commands[name]; !ok {
			t.Errorf("shortDesc has %q which is not in commands map", name)
		}
	}
}

func TestVersionDefault(t *testing.T) {
	if Version != "dev" {
		t.Errorf("default Version = %q, want %q", Version, "dev")
	}
}

func TestHelpContainsAllCommands(t *testing.T) {
	withArgs(t, []string{"sg"})
	out := captureStdout(t, func() {
		_ = Execute()
	})
	for _, name := range commandOrder {
		if !strings.Contains(out, name) {
			t.Errorf("help output missing command %q", name)
		}
	}
}
