package commands

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/mcbattirola/gitnotes/pkg/errflags"
	"github.com/mcbattirola/gitnotes/pkg/gn"
)

func Push(app *gn.GN, args []string) int {
	// gn push
	// pushes notes changes to origin
	pushCmd := flag.NewFlagSet("push", flag.ExitOnError)
	pushCmd.StringVar(&app.RemoteURL, "u", app.RemoteURL, "url of the remote origin in case none is set in the repository")
	pushCmd.Usage = func() {
		fmt.Println("Push notes changes to remote. Commits any uncommited change.")
		pushCmd.PrintDefaults()
	}

	if err := pushCmd.Parse(args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing command flags: %s\n", err.Error())
	}

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

	return 0
}
