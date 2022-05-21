package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const gnDirPath = "/home/matheus/Documents/projects/gitnotes" // TODO read this from config file

func main() {
	// TODO split args, remove bin name
	run(os.Args)
}

func run(args []string) int {
	// todo use proper flags
	// TODO is init necessary? maybe we can just do everything it does on edit
	switch args[1] {
	case "init":
		{
			// check if local repo exists (in config file or default path)
			//  create it if dont exist
			// create a dir to the current project on the notes repo if it doesnt exist (warn if it does)
		}
	case "edit":
		{
			// TODO check current branch
			// open file in the path (notesdir(from configfile)/current_repo/current_branch)
			_, err := os.Stat(gnDirPath)
			if os.IsNotExist(err) {
				err := os.MkdirAll(gnDirPath, os.ModeDir|0700)
				if err != nil {
					if errors.Is(err, fs.ErrPermission) {
						fmt.Fprintf(os.Stderr, "No permision to create directory %s", gnDirPath)
					}
					return 1
				}
			} else if err != nil {
				// TODO improve error handling
				panic(err)
			}

			// read current project name and branch
			project, err := getProjectRoot()
			if err != nil {
				panic(err)
			}

			// get user working repo
			dir, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			r, err := git.PlainOpen(dir)
			if err != nil {
				panic(err)
			}

			branch := getCurrentBranch(r)
			fmt.Println(branch)

			projectPath := fmt.Sprintf("%s/%s", gnDirPath, project)	
			notesPath := fmt.Sprintf("%s/%s", projectPath, branch)

			// we make the directory with the notespath instead of project path because
			// a branch name may contain slashes. In that case, we want to make the full path
			// and the slashes in branch name will become directories (which is ok for now)
			if err := os.MkdirAll(filepath.Dir(notesPath), os.ModeDir|0700); err != nil {
				panic(err)
			}

			_, err = os.OpenFile(notesPath, os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println(err.Error())
				panic(err.Error())
			}

			cmd := exec.Command("nvim", notesPath)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				panic(err.Error())
			}
		}
	case "sync":
		{
			// check in the config file if remote exists
			//  if it does, commit to it
			//  else, create it (will ask gh credentials or key or w/e)
		}
	}

	return 0
}

// getProjectRoot runs git through a syscall to get the top level directory
// we do it that way because go-git does not implement rev-parse
func getProjectRoot() (string, error) {
	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}

	s := strings.Split(strings.TrimSpace(string(path)), "/")
	return s[len(s)-1], nil
}

func getCurrentBranch(r *git.Repository) string {
	h, err := r.Reference(plumbing.HEAD, false)
	if err != nil {
		panic(err)
	}

	target := h.Target()

	s := strings.Split(string(target), "refs/heads/")
	if len(s) < 2 {
		return ""
	}
	return s[1]
}
