package main

import (
	"os"

	"github.com/mcbattirola/gitnotes/pkg/cli"
)

func main() {
	os.Exit(cli.Run(os.Args))
}
