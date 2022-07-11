package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/mcbattirola/gitnotes/pkg/gn"
)

// Commit commits the gitnotes
func Commit(app *gn.GN, args []string) int {
	// gn commit
	commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)
	commitCmd.StringVar(&app.CommitMessage, "message", app.CommitMessage, "commit message, in quotes")
	commitCmd.Usage = func() {
		fmt.Println("Commits notes changes. Example: gn commit --message \"Update notes\"")
		commitCmd.PrintDefaults()
	}

	if err := commitCmd.Parse(args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing commit command arguments: %s", err.Error())
		return 1
	}
	if err := app.Commit(); err != nil {
		fmt.Fprintf(os.Stderr, "error while commiting notes: %s\n", err.Error())
		return 1
	}

	return 0
}
