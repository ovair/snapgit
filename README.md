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

**Build from source** (requires Go 1.24+):

```bash
git clone https://github.com/ovair/snapgit.git
cd snapgit
go build -o sg ./cmd/sg
```

Move the `sg` binary somewhere in your `PATH`:

```bash
# Linux / macOS
sudo mv sg /usr/local/bin/

# Or add to your local bin
mv sg ~/.local/bin/
```

## Commands

| Command | What it does | Git equivalent |
|---|---|---|
| `sg create` | Create a new repo | `git init` |
| `sg get <url>` | Clone a repo | `git clone <url>` |
| `sg status` | Show repo status | `git status` |
| `sg add <file\|.>` | Stage file(s) | `git add <file>` |
| `sg save "msg"` | Stage all + commit | `git add -A && git commit -m "msg"` |
| `sg diff` | Show changes | `git diff` |
| `sg log` | Show history | `git log --oneline --graph --decorate` |
| `sg branch` | List branches | `git branch` |
| `sg new <branch>` | Create + switch branch | `git switch -c <branch>` |
| `sg go <branch>` | Switch branch | `git switch <branch>` |
| `sg fetch` | Fetch remote updates | `git fetch` |
| `sg pull` | Pull remote changes | `git pull` |
| `sg send` | Push to remote | `git push` |

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
```

## Philosophy

- **Human intention over Git mechanics** — commands describe what you want (`save`, `send`, `go`), not Git internals (`checkout`, `reset`, `rebase`)
- **Simplicity** — small, memorable command set covering daily workflows
- **Safety** — dangerous operations are excluded until proper safeguards exist
- **Compatibility** — works with existing Git repos, GitHub, and all remotes

## Requirements

- Git installed and available in `PATH`
- Go 1.24+ (build only)

## License

MIT
