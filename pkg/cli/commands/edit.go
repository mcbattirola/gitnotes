package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/mcbattirola/gitnotes/pkg/errflags"
	"github.com/mcbattirola/gitnotes/pkg/gn"
)

func Edit(app *gn.GN, args []string) int {
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

	if err := editCmd.Parse(args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing parameters: %s\n", err.Error())
		return 1
	}
	if err := checkEditParams(app); err != nil {
		fmt.Fprintf(os.Stderr, "error validating parameters: %s\n", err.Error())
		return 1
	}

	if err := app.Edit(); err != nil {
		fmt.Fprintf(os.Stderr, "error while editing file: %s\n", err.Error())
		return 1
	}

	return 0
}

func checkEditParams(app *gn.GN) error {
	if app.Project != "" && app.Branch == "" {
		return errflags.New("branch is necessary when specifying a project", errflags.BadParameter)
	}

	return nil
}
