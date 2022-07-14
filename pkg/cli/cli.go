// package cli reads flags, populate a struct of dependencies and run gn
package cli

import (
	"flag"
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
	// TODO this doesnt work
	var debug bool
	flag.BoolVar(&debug, "d", false, "enable debug logs")
	flag.Parse()

	app := gn.New(debug)

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
		"edit": {
			exec: commands.Edit,
			help: "edit the git notes",
		},
		"push": {
			exec: commands.Push,
			help: "push notes to remote",
		},
		"pull": {
			exec: commands.Pull,
			help: "pull notes from remote",
		},
		"commit": {
			exec: commands.Commit,
			help: "commit notes",
		},
		"path": {
			exec: commands.Path,
			help: "prints the notes path to stdio",
		},
		"delete": {
			exec: commands.Delete,
			help: "delete notes",
		},
	}

	subcommandIndex := 1
	if debug {
		subcommandIndex = 2
	}

	if len(os.Args) < subcommandIndex+1 {
		printSubcommansdHelp(cmds)
		return 1
	}

	cmd, ok := cmds[args[subcommandIndex]]
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
