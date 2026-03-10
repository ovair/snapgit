# SnapGit Development Plan

## Project Context

SnapGit (`sg`) is a human-friendly Git CLI wrapper written in Go (1.24+). It translates simple, intention-based commands into Git operations. Zero external dependencies — pure stdlib.

**Repo:** `github.com/ovair/snapgit`
**Current version:** v0.3.0
**Binary:** `sg`
**Module:** `snapgit`

## Architecture

```
cmd/sg/main.go              → Entry point (15 lines)
internal/cli/root.go        → Command dispatcher, auto-generated help, version
internal/cli/<command>.go   → One file per command (27 commands)
internal/cli/completions.go → Shell completion generators (bash/zsh/fish/powershell)
internal/cli/root_test.go   → CLI tests (~850 lines, mocked git.Run/RunOutput)
internal/git/runner.go      → Git execution: path caching, signal forwarding, exit codes, RunOutput
internal/git/runner_test.go → Git layer tests (130 lines)
Makefile                    → build, test, vet, check, install, clean
.goreleaser.yml             → Multi-platform releases + Homebrew tap (linux/mac/windows, amd64/arm64)
.github/workflows/ci.yml   → CI: vet + test + build on push/PR
.github/workflows/release.yml → Release on tag push via goreleaser
install.sh / install.ps1   → One-line install scripts
```

**Key design patterns:**

- `commands` map in root.go maps names → `{handler, usage, help}`
- `commandOrder` slice controls help display order
- `shortDesc` map holds one-line descriptions for help listing
- `printHelp()` auto-generates from `commandOrder` + `shortDesc` (no hardcoded help)
- `git.Run` / `git.RunOutput` / `git.RunExternal` are swappable function vars for testability
- Shared `execCmd()` helper for signal forwarding + exit code preservation
- Git binary path is resolved once and cached via `sync.Once`
- Consistency tests enforce `commands` ↔ `commandOrder` ↔ `shortDesc` sync

## What Exists (27 Commands)

| Command                  | Handler        | Git equivalent                            |
| ------------------------ | -------------- | ----------------------------------------- |
| `sg create`              | runCreate      | `git init`                                |
| `sg get <url>`           | runGet         | `git clone <url>`                         |
| `sg status`              | runStatus      | `git status`                              |
| `sg add <file\|.>`       | runAdd         | `git add <file>`                          |
| `sg save "msg"`          | runSave        | `git add -A && git commit -m "msg"`       |
| `sg diff`                | runDiff        | `git diff`                                |
| `sg log`                 | runLog         | `git log --oneline --graph --decorate`    |
| `sg branch`              | runBranch      | `git branch`                              |
| `sg new <branch>`        | runNew         | `git switch -c <branch>`                  |
| `sg go <branch>`         | runGo          | `git switch <branch>`                     |
| `sg fetch`               | runFetch       | `git fetch`                               |
| `sg pull`                | runPull        | `git pull`                                |
| `sg send`                | runSend        | `git push -u origin HEAD`                 |
| `sg undo`                | runUndo        | `git reset --soft HEAD~1`                 |
| `sg stash`               | runStash       | `git stash`                               |
| `sg pop`                 | runPop         | `git stash pop`                           |
| `sg merge <branch>`      | runMerge       | `git merge <branch>`                      |
| `sg tag [name]`          | runTag         | `git tag [name]`                          |
| `sg pr ["title"]`        | runPR          | `git push -u origin HEAD && gh pr create` |
| `sg rename <name>`       | runRename      | `git branch -m <name>`                    |
| `sg delete <branch>`     | runDelete      | `git branch -d <branch>`                  |
| `sg ignore <pattern>`    | runIgnore      | Appends to .gitignore                     |
| `sg whoami`              | runWhoami      | `git config user.name / user.email`       |
| `sg remote`              | runRemote      | `git remote -v`                           |
| `sg amend ["msg"]`       | runAmend       | `git commit --amend`                      |
| `sg release <version>`   | runRelease     | `git tag <v> && git push origin <v>`      |
| `sg completions <shell>` | runCompletions | Outputs shell completion script           |

Also: `sg help [command]`, `sg version` (aliases: --help/-h, --version/-v)

---

## v0.3.0 — COMPLETED

All items below were shipped in v0.3.0:

- [x] Fix `shortDesc` missing `"pr"` entry
- [x] Refactor `printHelp()` to auto-generate from maps
- [x] Refactor `runner.go`: deduplicate into `execCmd()`, add `RunOutput`
- [x] New commands: `rename`, `delete`, `ignore`, `whoami`, `remote`, `amend`
- [x] Shell completions: `sg completions <bash|zsh|fish|powershell>`
- [x] `sg release <version>` — tag + push in one command
- [x] Homebrew tap config in `.goreleaser.yml`
- [x] Consistency tests: `TestShortDescMatchesCommands`, `TestHelpContainsAllCommands`
- [x] README fully updated with all 27 commands + completions + Homebrew

---

## v0.4.0 — Bug Fixes & Hardening

### 1. Fix `ignore.go` — Windows newline bug

**Priority: HIGH**

- Line 21: `root = root[:len(root)-1]` assumes `\n` but Windows git outputs `\r\n`
- Fix: replace with `strings.TrimSpace(root)`

### 2. Fix `pr.go` — Check `gh` before pushing

**Priority: HIGH**

- Currently pushes the branch THEN tries `gh pr create`
- If `gh` is not installed, branch is pushed but PR fails — inconsistent state
- Fix: check `exec.LookPath("gh")` before running the push

### 3. Fix `delete.go` — Validate RunOutput result

**Priority: MEDIUM**

- If `git rev-parse --abbrev-ref HEAD` returns empty string, comparison fails silently
- Fix: check `strings.TrimSpace(current) == ""` and return an error

### 4. Fix `whoami.go` — Better error messages

**Priority: MEDIUM**

- If `user.name` fails, says "git user.name is not set" even if git itself isn't available
- Fix: check both configs, report which one(s) are missing

### 5. Fix `root.go` — Defensive label extraction

**Priority: LOW**

- Line ~121: `usage[len("sg "+name)+1:]` could panic if usage string is malformed
- Fix: add bounds check before slicing

### 6. Fix `runner.go` — Forward SIGTERM in addition to SIGINT

**Priority: MEDIUM**

- Only `os.Interrupt` (SIGINT) is forwarded to child processes
- Fix: add `syscall.SIGTERM` to `signal.Notify` (with build tag for non-Windows)

### 7. Fix install scripts — Checksum verification

**Priority: HIGH**

- Neither `install.sh` nor `install.ps1` verifies checksums
- goreleaser already generates `checksums.txt`
- Fix: download checksums.txt and verify before installing

### 8. Fix `install.sh` — POSIX compatibility

**Priority: HIGH**

- Uses `#!/bin/sh` but `${VERSION#v}` is not guaranteed POSIX
- Fix: either switch shebang to `#!/bin/bash` or use `echo "$VERSION" | sed 's/^v//'`

### 9. Fix `sg pr` command

**Priority: HIGH**

- Throws an error:
- PS C:\dev\Code\snapgit> sg pr
- branch 'claude/optimize-go-project-Ne3xv' set up to track 'origin/claude/optimize-go-project-Ne3xv'.
- Everything up-to-date
- To get started with GitHub CLI, please run: gh auth login
- Alternatively, populate the GH_TOKEN environment variable with a GitHub API authentication token.
- failed to create pull request: exit status 4

---

## v0.4.0 — New Features

### 9. `sg undo [n]` — Undo multiple commits

- Currently: `sg undo` → `git reset --soft HEAD~1`
- Enhancement: `sg undo 3` → `git reset --soft HEAD~3`
- Default to 1 if no argument
- Validate `n` is a positive integer

### 10. `sg diff --staged` — Show staged changes

- Currently: `sg diff` only shows unstaged changes
- Enhancement: `sg diff staged` → `git diff --cached`
- `sg diff` (no args) stays as-is: `git diff`

### 11. `sg stash "description"` — Named stashes

- Currently: `sg stash` → `git stash`
- Enhancement: `sg stash "work in progress"` → `git stash push -m "work in progress"`
- No args still does plain `git stash`

### 12. `sg log [n]` — Limit log output

- Currently: shows full log
- Enhancement: `sg log 10` → `git log --oneline --graph --decorate -10`
- No args shows all (current behavior)

### 13. `sg clean` — Remove untracked files (with confirmation)

- Equivalent: `git clean -fd`
- Must show what would be deleted first: `git clean -fdn` (dry run)
- Then ask for confirmation or require `--force` flag
- Safety-first approach

### 14. `sg revert <commit>` — Revert a specific commit

- Equivalent: `git revert <commit>`
- Requires 1 argument (commit hash or HEAD~n)
- Safe operation — creates a new commit that undoes changes

### 15. `sg cherry-pick <commit>` — Apply a commit from another branch

- Equivalent: `git cherry-pick <commit>`
- Requires 1 argument
- Useful for pulling individual fixes across branches

### 16. Create `ovair/homebrew-tap` repo

- The goreleaser Homebrew step fails because this repo doesn't exist yet
- Create empty public repo `ovair/homebrew-tap` on GitHub
- Re-release or manually push formula to verify
- Update install.sh to mention `brew install ovair/tap/sg`

---

## v0.4.0 — Optimizations

### 17. CI improvements

- Add `go fmt` check (fail if code isn't formatted)
- Add `go mod tidy` check (fail if go.sum is dirty)
- Add `govulncheck` for vulnerability scanning
- Pin Go version to exact match with `go.mod`

### 18. Standardize error messages

- All "usage" errors: `fmt.Errorf("usage: sg <command> <args>")`
- All git failures: `fmt.Errorf("failed to <action>: %w", err)`
- All missing tools: `fmt.Errorf("<tool> is required but not installed")`
- Audit all commands for consistency

### 19. `add` command — Support multiple files

- Currently: `sg add file.go` stages one file
- Enhancement: `sg add file1.go file2.go` → `git add file1.go file2.go`
- Pass `os.Args[2:]` instead of just `os.Args[2]`

### 20. `get` command — Support extra args

- Currently: `sg get <url>` only
- Enhancement: `sg get <url> <dir>` → `git clone <url> <dir>`
- Useful for cloning into a specific directory

### 21. Release version validation

- `sg release` accepts any string — `sg release oops` creates tag "oops"
- Validate version matches `v\d+\.\d+\.\d+` pattern
- Error with: `"version must be in format vX.Y.Z"`

### 22. Completion generators — Reduce duplication

- `branchCommands` map in completions.go is hardcoded
- All 4 shell generators duplicate branch-completing logic
- Refactor: single source of truth for which commands complete with branches

---

## v0.5.0 — Future Ideas

### 23. `sg bisect` — Binary search for bugs

- Wrapper around `git bisect start/good/bad/reset`
- Simplified interface: `sg bisect start`, `sg bisect good`, `sg bisect bad`, `sg bisect stop`

### 24. `sg blame <file>` — Show who changed what

- Equivalent: `git blame <file>`
- Simple passthrough

### 25. `sg conflicts` — Show merge conflict files

- Equivalent: `git diff --name-only --diff-filter=U`
- Useful during merge resolution

### 26. `sg upstream <remote-branch>` — Set tracking branch

- Equivalent: `git branch --set-upstream-to=origin/<branch>`
- Useful after creating branches without tracking

### 27. `--dry-run` global flag

- Shows what git commands would be run without executing
- Useful for learning and debugging
- Requires intercepting `git.Run` calls

### 28. Colored output

- Add color to `sg status`, `sg branch`, `sg log` output
- Use ANSI codes (stdlib only, no dependencies)
- Respect `NO_COLOR` environment variable
- Auto-detect TTY for piping support

### 29. `sg config` — Configuration file

- `~/.config/sg/config.toml` or similar
- Options: default branch name, auto-push after save, custom aliases
- Keep it simple — stdlib TOML/INI parsing

### 30. Man pages

- Generate man pages from command help text
- Install via `sg completions man` or build step
- Include in Homebrew formula

---

## Implementation Order (v0.4.0)

1. Bug fixes first (#1-#8) — all critical and medium priority
2. CI improvements (#17) — prevents future regressions
3. Standardize error messages (#18) — consistency pass
4. New features: `undo [n]`, `diff staged`, `stash "msg"`, `log [n]` (#9-#12)
5. New commands: `clean`, `revert`, `cherry-pick` (#13-#15)
6. Optimizations: multi-file add, get with dir, version validation (#19-#22)
7. Create homebrew-tap repo (#16)
8. Tests for everything new
9. Final README + PLAN update
10. `sg release v0.4.0`

## Rules

- Zero external dependencies — stdlib only
- One file per command in `internal/cli/`
- Every command gets: handler, entry in `commands` map, `commandOrder`, `shortDesc`, tests
- Tests mock `git.Run` / `git.RunOutput` — never call real git in CLI tests
- Keep the binary small and fast
- Follow existing code style exactly (no docstrings on handler funcs, terse error messages)
- `make check` must pass before any release
