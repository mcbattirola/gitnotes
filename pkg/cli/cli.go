// package cli reads flags, populate a struct of dependencies and run gn
package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/mcbattirola/gitnotes/pkg/config"
	"github.com/mcbattirola/gitnotes/pkg/errflags"
	"github.com/mcbattirola/gitnotes/pkg/gn"
)

func Run(args []string) int {
	app := gn.New()

	// order of precedence is: CLI > config file
	// 1. apply config file
	// 2. apply CLI args on top of it

	// open config file
	homeDir, err := os.UserHomeDir()
	configPath := homeDir + "/.config/gitnotes"
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading user home dir: %s\n", err.Error())
		return 1
	}
	configFileName := "gn.conf"
	if err := config.ReadConfigFile(app, configPath, configFileName); err != nil {
		fmt.Fprintf(os.Stderr, "error reading config file: %s\n", err.Error())
		return 1
	}

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

	// gn push
	// pushes notes changes to origin
	pushCmd := flag.NewFlagSet("push", flag.ExitOnError)
	pushCmd.StringVar(&app.RemoteURL, "u", app.RemoteURL, "url of the remote origin in case none is set in the repository")
	pushCmd.Usage = func() {
		fmt.Println("Push notes changes to remote. Commits any uncommited change.")
		pushCmd.PrintDefaults()
	}

	// gn pull
	// pulls from origin
	pullCmd := flag.NewFlagSet("pull", flag.ExitOnError)
	pullCmd.StringVar(&app.RemoteURL, "u", app.RemoteURL, "url of the remote origin in case none is set in the repository")
	pullCmd.Usage = func() {
		fmt.Println("Pull notes from origin")
		pullCmd.PrintDefaults()
	}

	// gn path
	// prints the notes path into stdout
	pathCmd := flag.NewFlagSet("path", flag.ExitOnError)
	pathCmd.Usage = func() {
		fmt.Println("Prints the notes paths to stdout.")
	}

	// gn commit
	commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)
	commitCmd.StringVar(&app.CommitMessage, "message", app.CommitMessage, "commit message, in quotes")
	commitCmd.Usage = func() {
		fmt.Println("Commits notes changes. Example: gn commit --message \"Update notes\"")
		commitCmd.PrintDefaults()
	}

	// gn delete
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteCmd.StringVar(&app.Project, "p", app.Project, "project to delete notes")
	deleteCmd.StringVar(&app.Branch, "b", app.Branch, "branch to delete notes")
	deleteCmd.Usage = func() {
		fmt.Println("delete notes")
		deleteCmd.PrintDefaults()
	}

	if len(os.Args) < 2 {
		fmt.Println("subcommand missing") // TODO print help, print subcommands
		return 1
	}

	switch args[1] {
	case editCmd.Name():
		{
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
		}
	case pushCmd.Name():
		{
			pushCmd.Parse(args[2:])
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
		}
	case pullCmd.Name():
		{
			pullCmd.Parse(args[2:])

			if err := app.Pull(); err != nil {
				if errflags.HasFlag(err, errflags.NoRemote) {
					fmt.Fprintf(os.Stderr, "no remote found. Try running again with -u [url] flag to set a remote origin\n")
					return 1
				}
				fmt.Fprintf(os.Stderr, "error pulling notes: %s\n", err.Error())
				return 1
			}

		}
	case pathCmd.Name():
		{
			fmt.Println(app.Path())
		}
	case commitCmd.Name():
		{
			if err := commitCmd.Parse(args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "error parsing commit command arguments: %s", err.Error())
				return 1
			}
			if err := app.Commit(); err != nil {
				fmt.Fprintf(os.Stderr, "error while commiting notes: %s\n", err.Error())
				return 1
			}
		}
	case deleteCmd.Name():
		{
			if err := deleteCmd.Parse(args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "error parsing delete command arguments: %s", err.Error())
				return 1
			}
			if err := app.Delete(); err != nil {
				fmt.Fprintf(os.Stderr, "error while deleting note: %s\n", err.Error())
				return 1
			}
		}
	default:
		{
			fmt.Fprintf(os.Stderr, "command %s not found.\n", args[1])
		}
	}

	return 0
}

func checkEditParams(app *gn.GN) error {
	if app.Project != "" && app.Branch == "" {
		return errflags.New("branch is necessary when specifying a project", errflags.BadParameter)
	}

	return nil
}
