package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func runCompletions() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: sg completions <bash|zsh|fish|powershell>")
	}

	shell := os.Args[2]
	switch shell {
	case "bash":
		fmt.Print(bashCompletion())
	case "zsh":
		fmt.Print(zshCompletion())
	case "fish":
		fmt.Print(fishCompletion())
	case "powershell":
		fmt.Print(powershellCompletion())
	default:
		return fmt.Errorf("unsupported shell: %s (supported: bash, zsh, fish, powershell)", shell)
	}
	return nil
}

// commandNames returns all command names from commandOrder.
func commandNames() []string {
	names := make([]string, len(commandOrder))
	copy(names, commandOrder)
	return append(names, "help", "version", "completions")
}

// branchCommands are commands that complete with branch names.
var branchCommands = map[string]bool{
	"go": true, "merge": true, "delete": true,
}

// branchCmdPattern returns a shell-friendly pattern for branch-completing commands.
func branchCmdPattern(sep string) string {
	cmds := make([]string, 0, len(branchCommands))
	for cmd := range branchCommands {
		cmds = append(cmds, cmd)
	}
	sort.Strings(cmds)
	return strings.Join(cmds, sep)
}

func bashCompletion() string {
	cmds := strings.Join(commandNames(), " ")
	pattern := branchCmdPattern("|")

	return fmt.Sprintf(`_sg_completions() {
    local cur prev commands
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    commands="%s"

    if [[ ${COMP_CWORD} -eq 1 ]]; then
        COMPREPLY=($(compgen -W "${commands}" -- "${cur}"))
        return 0
    fi

    case "${prev}" in
        %s)
            local branches
            branches=$(git branch --format='%%(refname:short)' 2>/dev/null)
            COMPREPLY=($(compgen -W "${branches}" -- "${cur}"))
            return 0
            ;;
        help)
            COMPREPLY=($(compgen -W "${commands}" -- "${cur}"))
            return 0
            ;;
    esac
}

complete -F _sg_completions sg
`, cmds, pattern)
}

func zshCompletion() string {
	var cmdList strings.Builder
	for _, name := range commandOrder {
		desc := shortDesc[name]
		cmdList.WriteString(fmt.Sprintf("        '%s:%s'\n", name, desc))
	}
	cmdList.WriteString("        'help:Show help for a command'\n")
	cmdList.WriteString("        'version:Show version'\n")
	cmdList.WriteString("        'completions:Generate shell completions'\n")

	pattern := branchCmdPattern("|")

	return fmt.Sprintf(`#compdef sg

_sg() {
    local -a commands
    commands=(
%s    )

    _arguments -C \
        '1: :->command' \
        '*:: :->args'

    case $state in
        command)
            _describe -t commands 'sg command' commands
            ;;
        args)
            case $words[1] in
                %s)
                    local branches
                    branches=(${(f)"$(git branch --format='%%(refname:short)' 2>/dev/null)"})
                    _describe -t branches 'branch' branches
                    ;;
                help)
                    _describe -t commands 'sg command' commands
                    ;;
            esac
            ;;
    esac
}

_sg "$@"
`, cmdList.String(), pattern)
}

func fishCompletion() string {
	var b strings.Builder
	b.WriteString("# Fish completions for sg\n")
	b.WriteString("complete -c sg -e\n\n")

	for _, name := range commandOrder {
		desc := shortDesc[name]
		b.WriteString(fmt.Sprintf("complete -c sg -n '__fish_use_subcommand' -a '%s' -d '%s'\n", name, desc))
	}
	b.WriteString("complete -c sg -n '__fish_use_subcommand' -a 'help' -d 'Show help for a command'\n")
	b.WriteString("complete -c sg -n '__fish_use_subcommand' -a 'version' -d 'Show version'\n")
	b.WriteString("complete -c sg -n '__fish_use_subcommand' -a 'completions' -d 'Generate shell completions'\n")
	b.WriteString("\n# Branch completions\n")
	cmds := make([]string, 0, len(branchCommands))
	for cmd := range branchCommands {
		cmds = append(cmds, cmd)
	}
	sort.Strings(cmds)
	for _, cmd := range cmds {
		b.WriteString(fmt.Sprintf("complete -c sg -n '__fish_seen_subcommand_from %s' -a '(git branch --format=\"%%(refname:short)\" 2>/dev/null)'\n", cmd))
	}
	b.WriteString("\n# Help completes command names\n")
	allCmds := strings.Join(commandNames(), " ")
	b.WriteString(fmt.Sprintf("complete -c sg -n '__fish_seen_subcommand_from help' -a '%s'\n", allCmds))

	return b.String()
}

func powershellCompletion() string {
	allCmds := strings.Join(commandNames(), "', '")
	// Build branch commands array for PowerShell
	cmds := make([]string, 0, len(branchCommands))
	for cmd := range branchCommands {
		cmds = append(cmds, "'"+cmd+"'")
	}
	sort.Strings(cmds)
	psBranchCmds := strings.Join(cmds, ", ")

	return fmt.Sprintf(`Register-ArgumentCompleter -CommandName sg -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $commands = @('%s')
    $args = $commandAst.CommandElements

    if ($args.Count -eq 1) {
        $commands | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
        }
        return
    }

    $subcommand = $args[1].ToString()
    switch ($subcommand) {
        { $_ -in @(%s) } {
            git branch --format='%%(refname:short)' 2>$null | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        'help' {
            $commands | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
    }
}
`, allCmds, psBranchCmds)
}
