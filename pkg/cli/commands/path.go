package commands

import (
	"flag"
	"fmt"

	"github.com/mcbattirola/gitnotes/pkg/gn"
)

func Path(app *gn.GN, args []string) int {
	// gn path
	// prints the notes path into stdout
	pathCmd := flag.NewFlagSet("path", flag.ExitOnError)
	pathCmd.Usage = func() {
		fmt.Println("Prints the notes paths to stdout.")
	}

	fmt.Println(app.Path())

	return 0
}
