<p align="center">
  <img src="public/snapgit.png" alt="SnapGit" width="120" />
</p>

<h1 align="center">SnapGit</h1>

<p align="center">
  Human-friendly Git commands. Same repos, same remotes, simpler words.
</p>

---

SnapGit (`sg`) is a CLI wrapper around Git that translates simple, intention-based commands into real Git operations. It doesn't replace Git — it works alongside it in any existing repository.

## Install

**Homebrew (macOS / Linux):**

```bash
brew install ovair/tap/sg
```

**Linux / macOS (script):**

```bash
curl -fsSL https://raw.githubusercontent.com/ovair/snapgit/main/install.sh | bash
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/ovair/snapgit/main/install.ps1 | iex
```

**Build from source** (requires Go 1.24+):

```bash
git clone https://github.com/ovair/snapgit.git
cd snapgit
go build -o sg ./cmd/sg
```

## Commands

| Command | What it does | Git equivalent |
|---|---|---|
| `sg create` | Create a new repo | `git init` |
| `sg get <url> [dir]` | Clone a repo | `git clone <url> [dir]` |
| `sg status` | Show repo status | `git status` |
| `sg add <file...>` | Stage file(s) | `git add <file...>` |
| `sg save "msg"` | Stage all + commit | `git add -A && git commit -m "msg"` |
| `sg diff [staged]` | Show changes (or staged) | `git diff` / `git diff --cached` |
| `sg log [n]` | Show history (last n) | `git log --oneline --graph --decorate [-n]` |
| `sg branch` | List branches | `git branch` |
| `sg new <branch>` | Create + switch branch | `git switch -c <branch>` |
| `sg go <branch>` | Switch branch | `git switch <branch>` |
| `sg fetch` | Fetch remote updates | `git fetch` |
| `sg pull` | Pull remote changes | `git pull` |
| `sg send` | Push to remote | `git push -u origin HEAD` |
| `sg undo [n]` | Undo last n commits (keep changes) | `git reset --soft HEAD~n` |
| `sg stash ["msg"]` | Stash working changes | `git stash` / `git stash push -m "msg"` |
| `sg pop` | Restore stashed changes | `git stash pop` |
| `sg merge <branch>` | Merge a branch | `git merge <branch>` |
| `sg tag [name]` | List or create tags | `git tag [name]` |
| `sg pr ["title"]` | Push + create GitHub PR | `git push -u origin HEAD && gh pr create` |
| `sg rename <name>` | Rename current branch | `git branch -m <name>` |
| `sg delete <branch>` | Delete a local branch | `git branch -d <branch>` |
| `sg ignore <pattern>` | Add pattern to .gitignore | *(file append)* |
| `sg whoami` | Show git user config | `git config user.name / user.email` |
| `sg remote` | Show remote URLs | `git remote -v` |
| `sg amend ["msg"]` | Amend last commit | `git commit --amend` |
| `sg clean [--force]` | Remove untracked files | `git clean -fd` |
| `sg revert <commit>` | Revert a commit | `git revert <commit>` |
| `sg cherry-pick <commit>` | Apply commit from another branch | `git cherry-pick <commit>` |
| `sg release <version>` | Tag and push a version | `git tag <v> && git push origin <v>` |
| `sg completions <shell>` | Generate shell completions | — |

## Quick Example

```bash
# Start a new project
sg create

# Make changes, then save them
sg save "add user login page"

# Push to remote
sg send

# Create a feature branch
sg new feature/dark-mode

# Switch back
sg go main

# Show last 10 commits
sg log 10

# View staged changes
sg diff staged

# Stash with a description
sg stash "work in progress"

# Undo last 2 commits (keep changes)
sg undo 2

# Stage multiple files
sg add main.go utils.go

# Rename your branch
sg rename feature/dark-theme

# Amend the last commit message
sg amend "fix: correct login redirect"

# Revert a specific commit
sg revert abc123

# Clean untracked files
sg clean

# See who you are
sg whoami

# Add a pattern to .gitignore
sg ignore "*.log"

# Release a new version
sg release v1.0.0
```

## Shell Completions

Generate completion scripts for your shell:

**Bash:**

```bash
# Add to ~/.bashrc
eval "$(sg completions bash)"
```

**Zsh:**

```zsh
# Add to ~/.zshrc
eval "$(sg completions zsh)"
```

**Fish:**

```fish
# Add to ~/.config/fish/config.fish
sg completions fish | source
```

**PowerShell:**

```powershell
# Add to your PowerShell profile
sg completions powershell | Out-String | Invoke-Expression
```

Completions support:
- `sg <TAB>` — list all commands
- `sg go <TAB>` — list local branches
- `sg merge <TAB>` — list local branches
- `sg delete <TAB>` — list local branches
- `sg help <TAB>` — list all commands

## Philosophy

- **Human intention over Git mechanics** — commands describe what you want (`save`, `send`, `go`), not Git internals (`checkout`, `reset`, `rebase`)
- **Simplicity** — small, memorable command set covering daily workflows
- **Safety** — dangerous operations require confirmation or explicit flags
- **Compatibility** — works with existing Git repos, GitHub, and all remotes
- **Zero dependencies** — pure Go stdlib, no external packages

## Requirements

- Git installed and available in `PATH`
- Go 1.24+ (build only)
- `gh` CLI (only for `sg pr`)

## License

MIT
