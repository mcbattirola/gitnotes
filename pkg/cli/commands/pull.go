package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/mcbattirola/gitnotes/pkg/errflags"
	"github.com/mcbattirola/gitnotes/pkg/gn"
)

func Pull(app *gn.GN, args []string) int {
	pullCmd := flag.NewFlagSet("pull", flag.ExitOnError)
	pullCmd.StringVar(&app.RemoteURL, "u", app.RemoteURL, "url of the remote origin in case none is set in the repository")
	pullCmd.Usage = func() {
		fmt.Println("Pull notes from origin")
		pullCmd.PrintDefaults()
	}

	pullCmd.Parse(args[2:])

	if err := app.Pull(); err != nil {
		if errflags.HasFlag(err, errflags.NoRemote) {
			fmt.Fprintf(os.Stderr, "no remote found. Try running again with -u [url] flag to set a remote origin\n")
			return 1
		}
		fmt.Fprintf(os.Stderr, "error pulling notes: %s\n", err.Error())
		return 1
	}

	return 0
}
