// package cli reads flags, populate a struct of dependencies and run gn
package cli

import (
	"fmt"
	"os"

	"github.com/mcbattirola/gitnotes/pkg/cli/commands"
	"github.com/mcbattirola/gitnotes/pkg/config"
	"github.com/mcbattirola/gitnotes/pkg/gn"
)

// command is the interface each subcommand must implement
// in order to be called by CLI
type command struct {
	exec func(app *gn.GN, args []string) int
	help string
}

func Run(args []string) int {
	app := gn.New()

	// read config file
	homeDir, err := os.UserHomeDir()
	configPath := homeDir + "/.config/gitnotes"
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading user home dir: %s\n", err.Error())
		return 1
	}
	configFileName := "gn.conf"
	if err := config.ReadConfigFile(app, configPath, configFileName); err != nil {
		fmt.Fprintf(os.Stderr, "error reading config file: %s\n", err.Error())
		return 1
	}

	cmds := map[string]command{
		"edit": command{
			exec: commands.Edit,
			help: "edit the git notes",
		},
		"push": command{
			exec: commands.Push,
			help: "push notes to remote",
		},
		"pull": command{
			exec: commands.Pull,
			help: "pull notes from remote",
		},
		"commit": command{
			exec: commands.Commit,
			help: "commit notes",
		},
		"path": command{
			exec: commands.Path,
			help: "prints the notes path to stdio",
		},
		"delete": command{
			exec: commands.Delete,
			help: "delete notes",
		},
	}

	if len(os.Args) < 2 {
		printSubcommansdHelp(cmds)
		return 1
	}

	cmd, ok := cmds[args[1]]
	if !ok {
		printSubcommansdHelp(cmds)
		return 1
	}

	return cmd.exec(app, args)
}

func printSubcommansdHelp(cmds map[string]command) {
	fmt.Println("Available commands:")
	for key, val := range cmds {
		fmt.Printf("- %s: %s\n", key, val.help)
	}

	fmt.Println("run [command] -h for more details")
}
