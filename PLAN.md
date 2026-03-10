# SnapGit Development Plan

## Project Context

SnapGit (`sg`) is a human-friendly Git CLI wrapper written in Go (1.24+). It translates simple, intention-based commands into Git operations. Zero external dependencies — pure stdlib.

**Repo:** `github.com/ovair/snapgit`
**Current version:** v0.4.0
**Binary:** `sg`
**Module:** `snapgit`

## Architecture

```
cmd/sg/main.go              → Entry point (15 lines)
internal/cli/root.go        → Command dispatcher, auto-generated help, version
internal/cli/<command>.go   → One file per command (30 commands)
internal/cli/completions.go → Shell completion generators (bash/zsh/fish/powershell)
internal/cli/root_test.go   → CLI tests (~850 lines, mocked git.Run/RunOutput)
internal/git/runner.go      → Git execution: path caching, signal forwarding, exit codes, RunOutput
internal/git/signal_unix.go → SIGINT + SIGTERM forwarding (non-Windows)
internal/git/signal_windows.go → SIGINT forwarding (Windows)
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

## What Exists (30 Commands)

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
| `sg clean [--force]`     | runClean       | `git clean -fd` (with confirmation)       |
| `sg revert <commit>`     | runRevert      | `git revert <commit>`                     |
| `sg cherry-pick <commit>`| runCherryPick  | `git cherry-pick <commit>`                |
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

## v0.4.0 — COMPLETED

All items below were shipped in v0.4.0:

### Bug Fixes & Hardening

- [x] Fix `ignore.go` — Windows newline bug: replaced `root[:len(root)-1]` with `strings.TrimSpace(root)`
- [x] Fix `pr.go` — Check `gh` before pushing: added `exec.LookPath("gh")` check before push
- [x] Fix `delete.go` — Validate RunOutput result: added empty string check for current branch
- [x] Fix `whoami.go` — Better error messages: reports which config(s) are missing
- [x] Fix `root.go` — Defensive label extraction: added bounds check before slicing usage string
- [x] Fix `runner.go` — Forward SIGTERM: added `syscall.SIGTERM` via build-tagged files (unix/windows)
- [x] Fix install scripts — Checksum verification: both `install.sh` and `install.ps1` verify SHA256
- [x] Fix `install.sh` — POSIX compatibility: switched shebang to `#!/bin/bash`
- [x] Fix `sg pr` — Pre-flight `gh` check prevents pushing before auth fails

### New Features

- [x] `sg undo [n]` — Undo multiple commits with validation
- [x] `sg diff staged` — Show staged changes via `git diff --cached`
- [x] `sg stash "message"` — Named stashes via `git stash push -m`
- [x] `sg log [n]` — Limit log output with validation
- [x] `sg clean [--force]` — Remove untracked files with confirmation/force flag
- [x] `sg revert <commit>` — Revert a specific commit
- [x] `sg cherry-pick <commit>` — Apply a commit from another branch

### Optimizations

- [x] CI improvements: added `go fmt` check, `go mod tidy` check, `govulncheck`, pinned Go version to `go.mod`
- [x] Standardized error messages across all commands
- [x] `sg add` supports multiple files: `sg add file1.go file2.go`
- [x] `sg get` supports directory arg: `sg get <url> <dir>`
- [x] `sg release` validates version format `vX.Y.Z`
- [x] Completion generators refactored: `branchCmdPattern()` generates shell patterns from `branchCommands` map
- [x] README fully updated with all 30 commands
- [x] Tests for all new commands and enhancements

### Not implemented (deferred)

- [ ] Create `ovair/homebrew-tap` repo (#16) — requires GitHub repo creation, not a code change

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

## Implementation Order (v0.5.0)

1. New commands from future ideas (#23-#26)
2. `--dry-run` global flag (#27)
3. Colored output (#28)
4. Configuration file (#29)
5. Man pages (#30)
6. Tests for everything new
7. Create homebrew-tap repo (#16) — deferred from v0.4.0
8. Final README + PLAN update
9. `sg release v0.5.0`

## Rules

- Zero external dependencies — stdlib only
- One file per command in `internal/cli/`
- Every command gets: handler, entry in `commands` map, `commandOrder`, `shortDesc`, tests
- Tests mock `git.Run` / `git.RunOutput` — never call real git in CLI tests
- Keep the binary small and fast
- Follow existing code style exactly (no docstrings on handler funcs, terse error messages)
- `make check` must pass before any release
