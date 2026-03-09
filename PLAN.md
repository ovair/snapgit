# SnapGit Development Plan

## Project Context

SnapGit (`sg`) is a human-friendly Git CLI wrapper written in Go (1.24+). It translates simple, intention-based commands into Git operations. Zero external dependencies — pure stdlib.

**Repo:** `github.com/ovair/snapgit`
**Current version:** v0.2.0 (tag), next release: v0.3.0
**Binary:** `sg`
**Module:** `snapgit`

## Architecture

```
cmd/sg/main.go              → Entry point (15 lines)
internal/cli/root.go        → Command dispatcher, help system, version
internal/cli/<command>.go   → One file per command (19 commands)
internal/cli/root_test.go   → CLI tests (572 lines, mocked git.Run)
internal/git/runner.go      → Git execution: path caching, signal forwarding, exit codes
internal/git/runner_test.go → Git layer tests (130 lines)
Makefile                    → build, test, vet, check, install, clean
.goreleaser.yml             → Multi-platform releases (linux/mac/windows, amd64/arm64)
.github/workflows/ci.yml   → CI: vet + test + build on push/PR
.github/workflows/release.yml → Release on tag push via goreleaser
install.sh / install.ps1   → One-line install scripts
```

**Key design patterns:**
- `commands` map in root.go maps names → `{handler, usage, help}`
- `commandOrder` slice controls help display order
- `shortDesc` map holds one-line descriptions for help listing
- `git.Run` / `git.RunExternal` are swappable function vars for testability
- Git binary path is resolved once and cached via `sync.Once`
- Signal forwarding (SIGINT/SIGTERM) to child processes
- Exit codes preserved through `ExitError` type

## What Exists (19 Commands)

| Command | Handler | Git equivalent |
|---|---|---|
| `sg create` | runCreate | `git init` |
| `sg get <url>` | runGet | `git clone <url>` |
| `sg status` | runStatus | `git status` |
| `sg add <file\|.>` | runAdd | `git add <file>` |
| `sg save "msg"` | runSave | `git add -A && git commit -m "msg"` |
| `sg diff` | runDiff | `git diff` |
| `sg log` | runLog | `git log --oneline --graph --decorate` |
| `sg branch` | runBranch | `git branch` |
| `sg new <branch>` | runNew | `git switch -c <branch>` |
| `sg go <branch>` | runGo | `git switch <branch>` |
| `sg fetch` | runFetch | `git fetch` |
| `sg pull` | runPull | `git pull` |
| `sg send` | runSend | `git push -u origin HEAD` |
| `sg undo` | runUndo | `git reset --soft HEAD~1` |
| `sg stash` | runStash | `git stash` |
| `sg pop` | runPop | `git stash pop` |
| `sg merge <branch>` | runMerge | `git merge <branch>` |
| `sg tag [name]` | runTag | `git tag [name]` |
| `sg pr ["title"]` | runPR | `git push -u origin HEAD && gh pr create` |

Also: `sg help [command]`, `sg version` (aliases: --help/-h, --version/-v)

## What's New — Features to Build

### 1. Shell Completions

Generate completion scripts for bash, zsh, fish, and PowerShell.

**Implementation:**
- Add `internal/cli/completions.go` with completion generators
- New command: `sg completions <shell>` — outputs completion script to stdout
- Shells: `bash`, `zsh`, `fish`, `powershell`
- Complete command names + subcommand arguments where applicable (e.g., `sg go` completes branch names via `git branch --format='%(refname:short)'`)
- Add `"completions"` to `commands` map, `commandOrder`, `shortDesc`
- Add `completions_test.go` — test that each shell outputs non-empty valid script
- Document in README under a "Shell Completions" section

**Completion behavior:**
- `sg <TAB>` → list all commands
- `sg go <TAB>` → list local branches
- `sg merge <TAB>` → list local branches
- `sg get <TAB>` → no completion (URL)
- `sg help <TAB>` → list all commands

### 2. Homebrew Formula

Create a Homebrew tap for easy macOS/Linux installation.

**Implementation:**
- Create repo `github.com/ovair/homebrew-tap` (or add formula to this repo under `Formula/`)
- Formula file: `sg.rb` — download from GitHub releases, install binary
- Update `.goreleaser.yml` to auto-publish to the Homebrew tap on release:
  ```yaml
  brews:
    - repository:
        owner: ovair
        name: homebrew-tap
      homepage: https://github.com/ovair/snapgit
      description: "Human-friendly Git CLI"
      install: |
        bin.install "sg"
  ```
- Update README install section: `brew install ovair/tap/sg`
- Test locally with `brew install --build-from-source`

### 3. `sg rename <new-name>` — Rename Current Branch

**Implementation:**
- Add `internal/cli/rename.go`
- Equivalent: `git branch -m <new-name>`
- Requires 1 argument
- Add to `commands`, `commandOrder`, `shortDesc`
- Add tests in `root_test.go`
- Update README commands table

### 4. `sg delete <branch>` — Delete a Local Branch

**Implementation:**
- Add `internal/cli/delete.go`
- Equivalent: `git branch -d <branch>`
- Use `-d` (safe delete, not `-D`) — fails if unmerged, which is the safe default
- Prevent deleting current branch (check first, show error)
- Requires 1 argument
- Add to `commands`, `commandOrder`, `shortDesc`
- Add tests

### 5. `sg ignore <pattern>` — Add to .gitignore

**Implementation:**
- Add `internal/cli/ignore.go`
- Appends the pattern to `.gitignore` in repo root (create if missing)
- Requires 1 argument
- No git command needed — pure file I/O
- Add tests (use temp dir with .gitignore)

### 6. `sg whoami` — Show Git User Config

**Implementation:**
- Add `internal/cli/whoami.go`
- Runs `git config user.name` and `git config user.email`
- Prints: `Name <email>`
- Add to commands map, tests

### 7. `sg remote` — Show Remote URL

**Implementation:**
- Add `internal/cli/remote.go`
- Equivalent: `git remote -v`
- No arguments needed
- Add to commands map, tests

### 8. `sg amend` — Amend Last Commit Message

**Implementation:**
- Add `internal/cli/amend.go`
- `sg amend "new message"` → `git commit --amend -m "new message"`
- `sg amend` (no args) → `git commit --amend --no-edit` (add staged changes to last commit)
- Add to commands map, tests

## Optimizations

### 9. README Sync

The README commands table is out of date — missing: `undo`, `stash`, `pop`, `merge`, `tag`, `pr`. Update it to list all commands and add sections for:
- Shell completions install instructions
- Homebrew install
- New commands

### 10. `shortDesc` Map Completeness

`shortDesc` is missing the `"pr"` entry. Add it: `"pr": "Push and create a pull request"`. Also add entries for every new command.

### 11. Help Output Sync

The `printHelp()` function has a hardcoded help string. It should be generated from `commandOrder` + `shortDesc` to stay in sync automatically. Refactor `printHelp()` to iterate `commandOrder` and format from `shortDesc` + command usage, eliminating the hardcoded block.

### 12. Test Coverage Expansion

- Ensure every new command has tests in `root_test.go`
- Add edge case tests: running commands outside a git repo, empty args
- Add a test that validates every entry in `commands` has a corresponding `shortDesc` entry
- Add a test that validates `commandOrder` contains every key in `commands`

### 13. Version Bump

After all features land:
- Tag as `v0.3.0`
- Verify `sg version` shows correct tag
- Push tag to trigger release pipeline

## Implementation Order

1. Fix `shortDesc` missing `"pr"` entry (quick fix)
2. Refactor `printHelp()` to auto-generate from maps (prevents future drift)
3. Update README with all current commands
4. Add new commands: `rename`, `delete`, `ignore`, `whoami`, `remote`, `amend`
5. Shell completions (`sg completions <shell>`)
6. Homebrew formula (goreleaser config + tap repo)
7. Tests for everything new
8. Final README update with completions + Homebrew sections
9. Tag v0.3.0, push, verify release

## Rules

- Zero external dependencies — stdlib only
- One file per command in `internal/cli/`
- Every command gets: handler, entry in `commands` map, `commandOrder`, `shortDesc`, tests
- Tests mock `git.Run` — never call real git in CLI tests
- Keep the binary small and fast
- Follow existing code style exactly (no docstrings on handler funcs, terse error messages)
- `make check` must pass before any release
