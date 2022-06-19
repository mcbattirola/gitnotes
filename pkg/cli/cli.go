// package cli reads flags, populate a struct of dependencies and run gn
package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/mcbattirola/gitnotes/pkg/config"
	"github.com/mcbattirola/gitnotes/pkg/errflags"
	"github.com/mcbattirola/gitnotes/pkg/gn"
)

func Run(args []string) int {
	app := gn.GN{}

	// order of precedence is: CLI > config file
	// 1. apply config file
	// 2. apply CLI args on top of it

	// open config file
	homeDir, err := os.UserHomeDir()
	configPath := homeDir + "/.config/gitnotes"
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading user home dir: %s\n", err.Error())
		return 1
	}
	configFileName := "gn.conf"
	if err := config.ReadConfigFile(&app, configPath, configFileName); err != nil {
		fmt.Fprintf(os.Stderr, "error reading config file: %s\n", err.Error())
		return 1
	}

	// read cli params
	// gn edit
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	editCmd.StringVar(&app.Editor, "e", app.Editor, "text editor")
	editCmd.StringVar(&app.Project, "p", app.Project, "project to edit notes")
	editCmd.StringVar(&app.Branch, "b", app.Branch, "branch to edit notes")
	editCmd.Usage = func() {
		fmt.Println("edit notes")
		editCmd.PrintDefaults()
	}

	// gn init
	// inits the git repo
	// initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	// gn push
	// pushes notes changes to origin
	pushCmd := flag.NewFlagSet("push", flag.ExitOnError)

	// gn path
	// prints the notes path into stdout
	pathCmd := flag.NewFlagSet("path", flag.ExitOnError)
	pathCmd.Usage = func() {
		fmt.Println("prints the notes path to stdout")
	}

	if len(os.Args) < 2 {
		fmt.Println("subcommand missing") // TODO print help, print subcommands
		return 1
	}

	switch args[1] {
	case "init":
		{
			// check if local repo exists (in config file or default path)
			//  create it if dont exist
			// create a dir to the current project on the notes repo if it doesnt exist (warn if it does)
			// dir default name should be gitnotes, should accept CLI arg
			fmt.Println("init not implemented")
			return 1
		}
	case editCmd.Name():
		{
			editCmd.Parse(args[2:])
			if err := checkInitParams(app); err != nil {
				fmt.Fprintf(os.Stderr, "error validating parameters: %s\n", err.Error())
				return 1
			}

			if err := app.Edit(); err != nil {
				fmt.Fprintf(os.Stderr, "error while editing file: %s\n", err.Error())
				return 1
			}
		}
	case "sync":
		{
			// check in the config file if remote exists
			//  if it does, commit to it
			//  else, create it (will ask gh credentials or key or w/e)
			fmt.Println("sync not implemented")
			return 1
		}
	case pushCmd.Name():
		{
			pushCmd.Parse(args[2:])
			if err := app.Push(); err != nil {
				// if remote was not found, prompt user to add a remote
				if errflags.HasFlag(err, errflags.NoRemote) {
					fmt.Printf("Remote not found. Enter remote URL: ")
					reader := bufio.NewReader(os.Stdin)
					url, _ := reader.ReadString('\n')

					// add remote origin and push again
					if err := app.AddOrigin(url); err != nil {
						fmt.Fprintf(os.Stderr, "error adding origin: %s\n", err.Error())
						return 1
					}
					if err := app.Push(); err != nil {
						fmt.Fprintf(os.Stderr, "error pushing notes: %s. Make sure the repository exists at %s", err.Error(), url)
						return 1
					}
					return 0
				}

				fmt.Fprintf(os.Stderr, "error pushing notes: %s\n", err.Error())
				return 1
			}
		}
	case pathCmd.Name():
		{
			fmt.Println(app.Path())
		}
	default:
		{
			fmt.Fprintf(os.Stderr, "command %s not found.\n", args[1])
		}
	}

	return 0
}

func checkInitParams(app gn.GN) error {
	if app.Project != "" && app.Branch == "" {
		return errflags.New("branch is necessary when specifying a project", errflags.BadParameter)
	}

	return nil
}
