package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/mcbattirola/gitnotes/pkg/errflags"
	"github.com/mcbattirola/gitnotes/pkg/gn"
)

func Print(app *gn.GN, args []string) int {
	// gn print
	// prints the notes into stdout
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	printCmd.StringVar(&app.Project, "p", app.Project, "project to edit notes")
	printCmd.StringVar(&app.Branch, "b", app.Branch, "branch to edit notes")
	printCmd.Usage = func() {
		fmt.Println("Prints the notes paths to stdout.")
	}

	if err := printCmd.Parse(args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing parameters: %s\n", err.Error())
		return 1
	}

	if err := checkPrintParams(app); err != nil {
		fmt.Fprintf(os.Stderr, "error validating parameters: %s\n", err.Error())
		return 1
	}

	_, err := fmt.Println(app.ReadNote())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error printing note content %s", err.Error())
		return 1
	}
	return 0
}

func checkPrintParams(app *gn.GN) error {
	if app.Project != "" && app.Branch == "" {
		return errflags.New("branch is necessary when specifying a project", errflags.BadParameter)
	}

	return nil
}
