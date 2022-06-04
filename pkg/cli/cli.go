// package cli reads flags, populate a struct of dependencies and run gn
package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/mcbattirola/gitnotes/pkg/config"
	"github.com/mcbattirola/gitnotes/pkg/gn"
)

func Run(args []string) int {
	// order of precedence is: CLI > config file
	// 1. apply config file
	// 2. apply CLI args on top of it

	app := gn.GN{}

	// open config file
	homeDir, err := os.UserHomeDir()
	configPath := homeDir + "/.config/gitnotes"
	if err != nil {
		panic(err)
	}
	configFileName := "gn.conf"
	if err := config.ReadConfigFile(&app, configPath, configFileName); err != nil {
		panic(err)
	}

	// read cli params
	// gn edit
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	editCmd.StringVar(&app.Editor, "editor", app.Editor, "text editor")
	editCmd.StringVar(&app.Project, "project", app.Project, "project to edit notes")
	editCmd.StringVar(&app.Branch, "branch", app.Branch, "branch to edit notes")
	editCmd.Usage = func() {
		fmt.Println("edits the git note")
		editCmd.PrintDefaults()
	}

	// gn init
	// inits the git repo
	// initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("subcommand missing") // TODO print help
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
	case "edit":
		{
			editCmd.Parse(args[2:])
			if app.Project != "" && app.Branch == "" {
				fmt.Printf("--branch is necessary when specifying a project")
				return 1
			}

			err := app.Edit()
			if err != nil {
				fmt.Printf("err on edit: %s \n", err.Error())
				panic(err)
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
	}

	return 0
}
