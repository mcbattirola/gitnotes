// package cli reads flags, populate a struct of dependencies and run gn
package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/mcbattirola/gitnotes/pkg/gn"
)

const defaultDirPath = "/home/matheus/Documents/projects/gitnotes" // TODO read this from config file
const defaultEditor = "vi"

func Run(args []string) int {
	// order of precedence is:
	// CLI > config file > defaults
	// we will read this in reverse order and always apply the latter one:
	// 1. create env with defaults
	// 2. apply config file
	// 3. apply CLI args on

	// defaults
	app := gn.GN{
		NotesPath: defaultDirPath,
		Editor:    defaultEditor,
		Project:   "",
		Branch:    "",
	}

	// open config file

	// cli
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	editCmd.StringVar(&app.Editor, "editor", app.Editor, "text editor")
	editCmd.StringVar(&app.Project, "project", app.Project, "project to edit note")
	editCmd.StringVar(&app.Branch, "branch", app.Branch, "branch to edit note")

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
			fmt.Println("init not implemented")
			return 1
		}
	case "edit":
		{
			editCmd.Parse(args[2:])
			// todo method to validate inputs?
			if app.Project != "" && app.Branch == "" {
				fmt.Printf("--branch is necessary when specifying a different project")
				return 1
			}

			err := app.Edit()
			if err != nil {
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
