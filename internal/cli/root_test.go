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
	if err := runDiff(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "diff" {
		t.Errorf("expected [diff], got %v", m.calls)
	}
}

func TestRunLog(t *testing.T) {
	m := withMockGit(t)
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
	if m.calls[0][0] != "push" {
		t.Errorf("expected [push], got %v", m.calls)
	}
}

func TestRunUndo(t *testing.T) {
	m := withMockGit(t)
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

func TestRunStash(t *testing.T) {
	m := withMockGit(t)
	if err := runStash(); err != nil {
		t.Fatal(err)
	}
	if m.calls[0][0] != "stash" {
		t.Errorf("expected [stash], got %v", m.calls)
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

func TestRunPR_NoTitle(t *testing.T) {
	m := withMockGit(t)
	ext := withMockExternal(t)
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

func TestRunPR_PushFails(t *testing.T) {
	m := withMockGit(t)
	m.err = fmt.Errorf("push failed")
	withMockExternal(t)
	withArgs(t, []string{"sg", "pr"})

	err := runPR()
	if err == nil {
		t.Fatal("expected error when push fails")
	}
	if !strings.Contains(err.Error(), "failed to push") {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- All commands registered ---

func TestAllCommandsRegistered(t *testing.T) {
	expected := []string{
		"create", "get", "status", "add", "save", "diff", "log",
		"branch", "new", "go", "fetch", "pull", "send",
		"undo", "stash", "pop", "merge", "tag", "pr",
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

func TestVersionDefault(t *testing.T) {
	if Version != "dev" {
		t.Errorf("default Version = %q, want %q", Version, "dev")
	}
}
