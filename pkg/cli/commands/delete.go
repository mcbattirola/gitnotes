package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/mcbattirola/gitnotes/pkg/gn"
)

func Delete(app *gn.GN, args []string) int {
	// gn delete
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteCmd.StringVar(&app.Project, "p", app.Project, "project to delete notes")
	deleteCmd.StringVar(&app.Branch, "b", app.Branch, "branch to delete notes")
	deleteCmd.Usage = func() {
		fmt.Println("delete notes")
		deleteCmd.PrintDefaults()
	}

	if err := deleteCmd.Parse(args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing delete command arguments: %s", err.Error())
		return 1
	}
	if err := app.Delete(); err != nil {
		fmt.Fprintf(os.Stderr, "error while deleting note: %s\n", err.Error())
		return 1
	}

	return 0
}
